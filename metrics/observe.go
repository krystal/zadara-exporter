package metrics

import (
	"context"
	"fmt"
	"strconv"

	"github.com/krystal/zadara-exporter/config"
	"github.com/krystal/zadara-exporter/zadara/commandcenter/vpsaobjectstorage"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type (
	// ClientFunc is a function that returns a ZadaraClient.
	ClientFunc func(ctx context.Context, target *config.Target) ZadaraClient
)

func (sm *StorageMetrics) observePolicy(
	o metric.Observer,
	policy *vpsaobjectstorage.ZiosStoragePolicy,
	attrs metric.MeasurementOption,
) error {
	// Observe the percentage of drives added metric.
	drivesAdded, err := strconv.ParseFloat(policy.PercentageDrivesAdded, 64)
	if err != nil {
		return fmt.Errorf("error parsing drives added: %w", err)
	}

	o.ObserveFloat64(sm.PercentageDrivesAdded, drivesAdded, attrs)

	o.ObserveFloat64(sm.RingBalanceNormalPercentage, policy.RingBalance.NormalPercentage, attrs)
	o.ObserveFloat64(sm.RingBalanceDegradedPercentage, policy.RingBalance.DegradedPercentage, attrs)
	o.ObserveFloat64(sm.RingBalanceCriticalPercentage, policy.RingBalance.CriticalPercentage, attrs)
	o.ObserveInt64(sm.FreeStorage, policy.FreeCapacity, attrs)
	o.ObserveInt64(sm.UsedStorage, policy.UsedCapacity, attrs)
	o.ObserveFloat64(sm.HealthPercentage, policy.HealthPercentage, attrs)
	o.ObserveFloat64(sm.RebalancePercentage, policy.RebalancePercentage, attrs)
	o.ObserveInt64(sm.RingBalanceNormalCount, policy.RingBalance.NormalCount, attrs)
	o.ObserveInt64(sm.RingBalanceDegradedCount, policy.RingBalance.DegradedCount, attrs)
	o.ObserveInt64(sm.RingBalanceCriticalCount, policy.RingBalance.CriticalCount, attrs)

	return nil
}

func (sm *StorageMetrics) observeStores(
	ctx context.Context,
	o metric.Observer,
	target *config.Target,
	client ZadaraClient,
) error {
	// Retrieve the storage policies from the ZadaraClient.
	stores, err := client.GetAllStoragePolicies(ctx)
	if err != nil {
		return fmt.Errorf("error getting storage policies: %w", err)
	}

	// Define the cloud name attribute.
	cloudNameAttr := attribute.String("cloud_name", target.CloudName)
	targeNameAttr := attribute.String("name", target.Name)

	for _, ssc := range stores {
		store := ssc.Store
		policies := ssc.Policies
		storeAttr := attribute.String("store", store.Name+"@"+target.CloudName)
		storeNameAttr := attribute.String("store_name", store.Name)

		storeLevelAttrs := metric.WithAttributes(
			targeNameAttr,
			cloudNameAttr,
			storeNameAttr,
			storeAttr,
		)

		o.ObserveInt64(sm.AccountsCount, store.AccountsCount, storeLevelAttrs)
		o.ObserveInt64(sm.UsersCount, store.UsersCount, storeLevelAttrs)
		o.ObserveInt64(sm.ContainersCount, store.ContainersCount, storeLevelAttrs)
		o.ObserveInt64(sm.ObjectsCount, store.ObjectsCount, storeLevelAttrs)
		o.ObserveInt64(sm.DrivesCount, store.Drives, storeLevelAttrs)
		o.ObserveInt64(sm.Cache, store.Cache, storeLevelAttrs)

		// Iterate over each policy.
		for _, policy := range policies {
			// Define the policy level attributes.
			policyLevelAttrs := metric.WithAttributes(
				targeNameAttr,
				cloudNameAttr,
				storeNameAttr,
				storeAttr,
				attribute.String("policy_name", policy.Name),
			)

			if err := sm.observePolicy(o, policy, policyLevelAttrs); err != nil {
				return err
			}
		}
	}

	return nil
}

// StorageMetricsObserve returns a metric callback function that observes storage metrics for the given targets.
// It takes a slice of targets and a newclient function as parameters.
// The newclient function is used to create a new client for each target.
// The metric callback function iterates over the targets,
// creates a client for each target using the newclient function,
// and calls the observeStores function to observe the storage metrics for the target using the client.
// If any error occurs during the observation, it is returned.
// If all observations are successful, nil is returned.
func (sm *StorageMetrics) StorageMetricsObserve(targets []*config.Target, newclient ClientFunc) metric.Callback {
	// Define the metric callback function.
	return func(ctx context.Context, o metric.Observer) error {
		for _, target := range targets {
			client := newclient(ctx, target)
			if err := sm.observeStores(ctx, o, target, client); err != nil {
				return err
			}
		}

		return nil
	}
}
