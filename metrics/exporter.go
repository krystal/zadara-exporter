package metrics

import (
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

const (
	// DefaultNamespace is the default namespace for the Prometheus exporter.
	DefaultNamespace = "zadara"

	// DefaultPath is the default path for the Prometheus exporter.
	DefaultPath = "/metrics"
)

// SetupPrometheusExporter initialises and sets up the Prometheus exporter for metrics.
// It creates a new Prometheus exporter, sets it as the meter provider, and returns any error encountered.
func SetupPrometheusExporter(namespace string) error {
	if namespace == "" {
		namespace = DefaultNamespace
	}

	exporter, err := prometheus.New(
		prometheus.WithNamespace(namespace),
		prometheus.WithoutScopeInfo(),
		prometheus.WithoutTargetInfo(),
	)
	if err != nil {
		return fmt.Errorf("failed to create prometheus exporter: %w", err)
	}

	provider := metric.NewMeterProvider(metric.WithReader(exporter))

	otel.SetMeterProvider(provider)

	return nil
}
