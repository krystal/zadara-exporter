package metrics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/krystal/zadara-exporter/metrics"
	"github.com/krystal/zadara-exporter/zadara/commandcenter/vpsaobjectstorage"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
)

type (
	mockZadaraClient struct {
		mock.Mock
	}
	mockObserver struct {
		embedded.Observer
		mock.Mock
	}
)

func (m *mockZadaraClient) GetAllStoragePolicies(
	ctx context.Context,
	cloudName string,
) ([]*vpsaobjectstorage.ZiosStoragePolicy, error) {
	args := m.Called(ctx, cloudName)

	firstArg, ok := args.Get(0).([]*vpsaobjectstorage.ZiosStoragePolicy)
	if !ok {
		return nil, fmt.Errorf("error with arg: %w", args.Error(1))
	}

	return firstArg, nil
}

func (m *mockObserver) ObserveInt64(obsrv metric.Int64Observable, value int64, opts ...metric.ObserveOption) {
	m.Called(obsrv, value, opts)
}

func (m *mockObserver) ObserveFloat64(obsrv metric.Float64Observable, value float64, opts ...metric.ObserveOption) {
	m.Called(obsrv, value, opts)
}

func TestStorageMetricsObserve(t *testing.T) {
	t.Parallel()

	// Create a mock ZadaraClient.
	mockClient := new(mockZadaraClient)

	// Create a StorageMetrics instance.

	meter := otel.Meter("zadara")
	storageMetrics, err := metrics.NewStorageMetrics(meter)
	require.NoError(t, err)

	// Create a test observer.
	observer := new(mockObserver)

	// Set up expectations for the mock client.
	mockClient.On("GetAllStoragePolicies", mock.Anything, "cloudName").Return([]*vpsaobjectstorage.ZiosStoragePolicy{
		{Name: "policy1", FreeCapacity: 100, UsedCapacity: 50},
		{Name: "policy2", FreeCapacity: 200, UsedCapacity: 150},
	}, nil)

	observer.ExpectedCalls = []*mock.Call{
		{
			Method:    "ObserveInt64",
			Arguments: mock.Arguments{storageMetrics.FreeStorage, int64(100), mock.Anything},
		},
		{
			Method:    "ObserveInt64",
			Arguments: mock.Arguments{storageMetrics.UsedStorage, int64(50), mock.Anything},
		},
		{
			Method:    "ObserveInt64",
			Arguments: mock.Arguments{storageMetrics.FreeStorage, int64(200), mock.Anything},
		},
		{
			Method:    "ObserveInt64",
			Arguments: mock.Arguments{storageMetrics.UsedStorage, int64(150), mock.Anything},
		},
	}

	// Call the function being tested.
	err = storageMetrics.StorageMetricsObserve(mockClient)(context.Background(), observer)

	// Assert that there were no errors.
	require.NoError(t, err)

	// Assert that the mock client's method was called with the expected arguments.
	mockClient.AssertCalled(t, "GetAllStoragePolicies", mock.Anything, "cloudName")
}
