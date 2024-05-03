// Package metrics : provides the metrics for the zadara-exporter.
package metrics

import (
	"context"
	"fmt"

	"github.com/krystal/zadara-exporter/zadara/commandcenter/vpsaobjectstorage"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type (
	// StorageMetrics : provides the metrics for the storage.
	StorageMetrics struct {
		FreeStorage metric.Int64ObservableGauge
		UsedStorage metric.Int64ObservableGauge
	}

	// ZadaraClient : provides the client for the Zadara storage.
	ZadaraClient interface {
		GetAllStoragePolicies(
			ctx context.Context,
			cloudName string,
		) ([]*vpsaobjectstorage.ZiosStoragePolicy, error)
	}
)

// NewStorageMetrics : creates a new storage metrics for the zadara-exporter.
func NewStorageMetrics(meter metric.Meter) (*StorageMetrics, error) {
	freeStorage, err := meter.Int64ObservableGauge("zadara_free_storage",
		metric.WithDescription("The amount of free storage in the Zadara storage."))
	if err != nil {
		return nil, fmt.Errorf("failed to create free storage gauge: %w", err)
	}

	usedStorage, err := meter.Int64ObservableGauge("zadara_used_storage",
		metric.WithDescription("The amount of used storage in the Zadara storage."))
	if err != nil {
		return nil, fmt.Errorf("failed to create used storage gauge: %w", err)
	}

	return &StorageMetrics{
		FreeStorage: freeStorage,
		UsedStorage: usedStorage,
	}, nil
}

// RegisterStorageMetrics : registers the storage metrics for the zadara-exporter.
func RegisterStorageMetrics(client ZadaraClient) error {
	meter := otel.Meter("zadara")

	metrics, err := NewStorageMetrics(meter)
	if err != nil {
		return fmt.Errorf("failed to create storage metrics: %w", err)
	}

	_, err = meter.RegisterCallback(metrics.StorageMetricsObserve(client),
		metrics.FreeStorage,
		metrics.UsedStorage,
	)
	if err != nil {
		return fmt.Errorf("failed to register storage metrics: %w", err)
	}

	return nil
}
