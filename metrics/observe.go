package metrics

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// StorageMetricsObserve : provides the metrics observable callback for the storage.
func (sm *StorageMetrics) StorageMetricsObserve(client ZadaraClient) metric.Callback {
	return func(ctx context.Context, o metric.Observer) error {
		policies, err := client.GetAllStoragePolicies(ctx, "cloudName")
		if err != nil {
			return fmt.Errorf("error getting storage policies: %w", err)
		}

		cloudNameAttr := attribute.String("cloud_name", "cloudName")

		for _, policy := range policies {
			policyNameAttr := attribute.String("policy_name", policy.Name)

			o.ObserveInt64(sm.FreeStorage, policy.FreeCapacity,
				metric.WithAttributes(
					cloudNameAttr,
					policyNameAttr,
				),
			)

			o.ObserveInt64(sm.UsedStorage, policy.UsedCapacity,
				metric.WithAttributes(
					cloudNameAttr,
					policyNameAttr,
				),
			)
		}

		return nil
	}
}
