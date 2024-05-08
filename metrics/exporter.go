package metrics

import (
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

// SetupPrometheusExporter initialises and sets up the Prometheus exporter for metrics.
// It creates a new Prometheus exporter, sets it as the meter provider, and returns any error encountered.
func SetupPrometheusExporter() error {
	exporter, err := prometheus.New()
	if err != nil {
		return fmt.Errorf("failed to create prometheus exporter: %w", err)
	}

	provider := metric.NewMeterProvider(metric.WithReader(exporter))

	otel.SetMeterProvider(provider)

	return nil
}
