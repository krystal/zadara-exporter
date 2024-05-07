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

func (sm *StorageMetrics) observePolicies(
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
	cloudNameAttr := attribute.String("cloud_name", "cloudName")
	targeNameAttr := attribute.String("name", target.Name)

	for _, ssc := range stores {
		store := ssc.Store
		policies := ssc.Policies
		storeNameAttr := attribute.String("store_name", store.Name)

		storeLevelAttrs := metric.WithAttributes(
			targeNameAttr,
			cloudNameAttr,
			storeNameAttr,
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
				attribute.String("policy_name", policy.Name),
			)

			if err := sm.observePolicies(o, policy, policyLevelAttrs); err != nil {
				return err
			}
		}
	}

	return nil
}

// StorageMetricsObserve returns a metric callback function that observes storage metrics.
// It takes a ZadaraClient as a parameter and returns an error.
// The callback function uses the provided ZadaraClient to retrieve storage policies and observe storage metrics.
// It iterates over each policy and observes the free and used storage capacities using the provided metric.Observer.
// The observed metrics include the cloud name and policy name as attributes.
// If an error occurs while retrieving the storage policies, it is returned as an error.
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
