package metrics

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/krystal/zadara-exporter/zadara/commandcenter/vpsaobjectstorage"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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

func newStorageMetrics(meter metric.Meter) (*StorageMetrics, error) {
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

	metrics, err := newStorageMetrics(meter)
	if err != nil {
		return fmt.Errorf("failed to create storage metrics: %w", err)
	}

	_, err = meter.RegisterCallback(func(ctx context.Context, o metric.Observer) error {
		policies, err := client.GetAllStoragePolicies(ctx, "cloudName")
		if err != nil {
			return fmt.Errorf("error getting storage policies: %w", err)
		}

		cloudNameAttr := attribute.String("cloud_name", "cloudName")

		for _, policy := range policies {
			policyNameAttr := attribute.String("policy_name", policy.Name)

			o.ObserveInt64(metrics.FreeStorage, policy.FreeCapacity,
				metric.WithAttributes(
					cloudNameAttr,
					policyNameAttr,
				),
			)

			o.ObserveInt64(metrics.UsedStorage, policy.UsedCapacity,
				metric.WithAttributes(
					cloudNameAttr,
					policyNameAttr,
				),
			)
		}

		return nil
	},
		metrics.FreeStorage,
		metrics.UsedStorage,
	)
	if err != nil {
		return fmt.Errorf("failed to register storage metrics: %w", err)
	}

	return nil
}

// Serve : serves the metrics for the zadara-exporter.
func Serve(ctx context.Context) error {
	mux := http.NewServeMux()

	// Create a new HTTP handler for serving the metrics.
	mux.Handle("/metrics", promhttp.Handler())

	const ReadHeaderTimeout = 10 * time.Second

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: ReadHeaderTimeout,
	}

	// Start the HTTP server in a separate goroutine.
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("error starting HTTP server: %v", err)
		}
	}()

	// Wait for the context to be done.
	<-ctx.Done()

	// Shutdown the HTTP server gracefully.
	err := server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("error shutting down HTTP server: %w", err)
	}

	return nil
}
