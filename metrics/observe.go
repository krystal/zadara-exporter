package metrics

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// StorageMetricsObserve returns a metric callback function that observes storage metrics.
// It takes a ZadaraClient as a parameter and returns an error.
// The callback function uses the provided ZadaraClient to retrieve storage policies and observe storage metrics.
// It iterates over each policy and observes the free and used storage capacities using the provided metric.Observer.
// The observed metrics include the cloud name and policy name as attributes.
// If an error occurs while retrieving the storage policies, it is returned as an error.
func (sm *StorageMetrics) StorageMetricsObserve(client ZadaraClient) metric.Callback {
	// Define the metric callback function.
	return func(ctx context.Context, o metric.Observer) error {
		// Retrieve the storage policies from the ZadaraClient.
		policies, err := client.GetAllStoragePolicies(ctx, "cloudName")
		if err != nil {
			return fmt.Errorf("error getting storage policies: %w", err)
		}

		// Define the cloud name attribute.
		cloudNameAttr := attribute.String("cloud_name", "cloudName")

		// Iterate over each policy.
		for _, policy := range policies {
			// Define the policy name attribute.
			policyNameAttr := attribute.String("policy_name", policy.Name)

			// Observe the free storage capacity metric.
			o.ObserveInt64(sm.FreeStorage, policy.FreeCapacity,
				metric.WithAttributes(
					cloudNameAttr,
					policyNameAttr,
				),
			)

			// Observe the used storage capacity metric.
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
