// Package metrics provides the metrics for the zadara-exporter.
package metrics

import (
	"context"
	"fmt"

	"github.com/krystal/zadara-exporter/config"
	"github.com/krystal/zadara-exporter/zadara/commandcenter"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type (
	// StorageMetrics provides the metrics for the storage.
	StorageMetrics struct {
		FreeStorage                   metric.Int64ObservableGauge
		UsedStorage                   metric.Int64ObservableGauge
		AccountsCount                 metric.Int64ObservableGauge
		UsersCount                    metric.Int64ObservableGauge
		ContainersCount               metric.Int64ObservableGauge
		ObjectsCount                  metric.Int64ObservableGauge
		DrivesCount                   metric.Int64ObservableGauge
		Cache                         metric.Int64ObservableGauge
		HealthPercentage              metric.Float64ObservableGauge
		RebalancePercentage           metric.Float64ObservableGauge
		PercentageDrivesAdded         metric.Float64ObservableGauge
		RingBalanceNormalPercentage   metric.Float64ObservableGauge
		RingBalanceDegradedPercentage metric.Float64ObservableGauge
		RingBalanceCriticalPercentage metric.Float64ObservableGauge
		RingBalanceNormalCount        metric.Int64ObservableGauge
		RingBalanceDegradedCount      metric.Int64ObservableGauge
		RingBalanceCriticalCount      metric.Int64ObservableGauge
	}

	// ZadaraClient provides the client for the Zadara storage.
	ZadaraClient interface {
		GetAllStoragePolicies(ctx context.Context) ([]*commandcenter.StoreStoragePolicies, error)
	}
)

func storeMetrics(meter metric.Meter, storageMetrics *StorageMetrics) error {
	var err error

	storageMetrics.AccountsCount, err = meter.Int64ObservableGauge("zadara_accounts_count",
		metric.WithDescription("The number of accounts in the Zadara store."))
	if err != nil {
		return fmt.Errorf("failed to create accounts count gauge: %w", err)
	}

	storageMetrics.UsersCount, err = meter.Int64ObservableGauge("zadara_users_count",
		metric.WithDescription("The number of users in the Zadara store."))
	if err != nil {
		return fmt.Errorf("failed to create users count gauge: %w", err)
	}

	storageMetrics.ContainersCount, err = meter.Int64ObservableGauge("zadara_containers_count",
		metric.WithDescription("The number of containers in the Zadara store."))
	if err != nil {
		return fmt.Errorf("failed to create containers count gauge: %w", err)
	}

	storageMetrics.ObjectsCount, err = meter.Int64ObservableGauge("zadara_objects_count",
		metric.WithDescription("The number of objects in the Zadara store."))
	if err != nil {
		return fmt.Errorf("failed to create objects count gauge: %w", err)
	}

	storageMetrics.DrivesCount, err = meter.Int64ObservableGauge("zadara_drives_count",
		metric.WithDescription("The number of drives in the Zadara store."))
	if err != nil {
		return fmt.Errorf("failed to create drives count gauge: %w", err)
	}

	storageMetrics.Cache, err = meter.Int64ObservableGauge("zadara_cache",
		metric.WithDescription("The amount of cache in the Zadara store."))
	if err != nil {
		return fmt.Errorf("failed to create cache gauge: %w", err)
	}

	return nil
}

func storagePolicyMetrics(meter metric.Meter, storageMetrics *StorageMetrics) error {
	var err error

	storageMetrics.FreeStorage, err = meter.Int64ObservableGauge("zadara_free_storage",
		metric.WithDescription("The amount of free storage in the Zadara store storage policy."))
	if err != nil {
		return fmt.Errorf("failed to create free storage gauge: %w", err)
	}

	storageMetrics.UsedStorage, err = meter.Int64ObservableGauge("zadara_used_storage",
		metric.WithDescription("The amount of used storage in the Zadara store storage policy."))
	if err != nil {
		return fmt.Errorf("failed to create used storage gauge: %w", err)
	}

	storageMetrics.HealthPercentage, err = meter.Float64ObservableGauge("zadara_health_percentage",
		metric.WithDescription("The percentage of health in the Zadara store."))
	if err != nil {
		return fmt.Errorf("failed to create health percentage gauge: %w", err)
	}

	storageMetrics.RebalancePercentage, err = meter.Float64ObservableGauge("zadara_rebalance_percentage",
		metric.WithDescription("The percentage of rebalance in the Zadara store."))
	if err != nil {
		return fmt.Errorf("failed to create rebalance percentage gauge: %w", err)
	}

	storageMetrics.PercentageDrivesAdded, err = meter.Float64ObservableGauge("zadara_percentage_drives_added",
		metric.WithDescription("The percentage of drives added in the Zadara store."))
	if err != nil {
		return fmt.Errorf("failed to create percentage drives added gauge: %w", err)
	}

	return nil
}

func ringBalanceMetrics(meter metric.Meter, storageMetrics *StorageMetrics) error {
	var err error

	storageMetrics.RingBalanceNormalPercentage, err = meter.Float64ObservableGauge("zadara_ring_balance_normal_percentage",
		metric.WithDescription("The percentage of normal ring balance in the Zadara store."))
	if err != nil {
		return fmt.Errorf("failed to create ring balance normal percentage gauge: %w", err)
	}

	storageMetrics.RingBalanceDegradedPercentage, err = meter.Float64ObservableGauge(
		"zadara_ring_balance_Degraded_percentage",
		metric.WithDescription("The percentage of Degraded ring balance in the Zadara store."))
	if err != nil {
		return fmt.Errorf("failed to create ring balance Degraded percentage gauge: %w", err)
	}

	storageMetrics.RingBalanceCriticalPercentage, err = meter.Float64ObservableGauge(
		"zadara_ring_balance_critical_percentage",
		metric.WithDescription("The percentage of critical ring balance in the Zadara store."))
	if err != nil {
		return fmt.Errorf("failed to create ring balance critical percentage gauge: %w", err)
	}

	storageMetrics.RingBalanceNormalCount, err = meter.Int64ObservableGauge("zadara_ring_balance_normal_count",
		metric.WithDescription("The count of normal ring balance in the Zadara store."))
	if err != nil {
		return fmt.Errorf("failed to create ring balance normal count gauge: %w", err)
	}

	storageMetrics.RingBalanceDegradedCount, err = meter.Int64ObservableGauge("zadara_ring_balance_Degraded_count",
		metric.WithDescription("The count of Degraded ring balance in the Zadara store."))
	if err != nil {
		return fmt.Errorf("failed to create ring balance Degraded count gauge: %w", err)
	}

	storageMetrics.RingBalanceCriticalCount, err = meter.Int64ObservableGauge("zadara_ring_balance_critical_count",
		metric.WithDescription("The count of critical ring balance in the Zadara store."))
	if err != nil {
		return fmt.Errorf("failed to create ring balance critical count gauge: %w", err)
	}

	return nil
}

// NewStorageMetrics creates a new instance of StorageMetrics using the provided meter.
// It returns a pointer to the created StorageMetrics and an error, if any.
func NewStorageMetrics(meter metric.Meter) (*StorageMetrics, error) {
	storageMetrics := &StorageMetrics{}

	if err := storeMetrics(meter, storageMetrics); err != nil {
		return nil, err
	}

	if err := storagePolicyMetrics(meter, storageMetrics); err != nil {
		return nil, err
	}

	if err := ringBalanceMetrics(meter, storageMetrics); err != nil {
		return nil, err
	}

	return storageMetrics, nil
}

// RegisterStorageMetrics registers storage metrics for the given Zadara client.
// It creates storage metrics using the provided meter and registers the metrics
// callback to observe the storage metrics for the client.
// Returns an error if there was a failure in creating or registering the metrics.
func RegisterStorageMetrics(targets []*config.Target) error {
	meter := otel.Meter("zadara")

	metrics, err := NewStorageMetrics(meter)
	if err != nil {
		return fmt.Errorf("failed to create storage metrics: %w", err)
	}

	_, err = meter.RegisterCallback(metrics.StorageMetricsObserve(targets,
		func(_ context.Context, target *config.Target) ZadaraClient {
			return commandcenter.NewClientFromToken(target.APIBaseURL, target.Token, target.CloudName)
		}),
		metrics.FreeStorage,
		metrics.UsedStorage,
		metrics.AccountsCount,
		metrics.UsersCount,
		metrics.ContainersCount,
		metrics.ObjectsCount,
		metrics.DrivesCount,
		metrics.Cache,
		metrics.HealthPercentage,
		metrics.RebalancePercentage,
		metrics.PercentageDrivesAdded,
		metrics.RingBalanceNormalPercentage,
		metrics.RingBalanceDegradedPercentage,
		metrics.RingBalanceCriticalPercentage,
		metrics.RingBalanceNormalCount,
		metrics.RingBalanceDegradedCount,
		metrics.RingBalanceCriticalCount,
	)
	if err != nil {
		return fmt.Errorf("failed to register storage metrics: %w", err)
	}

	return nil
}
