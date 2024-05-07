package metrics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/krystal/zadara-exporter/config"
	"github.com/krystal/zadara-exporter/metrics"
	"github.com/krystal/zadara-exporter/zadara/commandcenter"
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
) ([]*commandcenter.StoreStoragePolicies, error) {
	args := m.Called(ctx, cloudName)

	firstArg, ok := args.Get(0).([]*commandcenter.StoreStoragePolicies)
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

//nolint:funlen // majority of the length is mocking
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
	mockClient.On("GetAllStoragePolicies", mock.Anything, "cloudName").Return([]*commandcenter.StoreStoragePolicies{
		{
			Store: &vpsaobjectstorage.Zios{
				AccountsCount:   2,
				UsersCount:      45,
				ContainersCount: 23,
				ObjectsCount:    78,
				Drives:          8,
				Cache:           1670,
			},
			Policies: []*vpsaobjectstorage.ZiosStoragePolicy{
				{
					Name:                  "policy1",
					FreeCapacity:          100,
					UsedCapacity:          50,
					HealthPercentage:      99.9,
					RebalancePercentage:   50.0,
					PercentageDrivesAdded: "55.2",
					RingBalance: vpsaobjectstorage.RingBalance{
						NormalPercentage:   75.0,
						DegradedPercentage: 12.5,
						CriticalPercentage: 12.5,
						NormalCount:        150,
						DegradedCount:      25,
						CriticalCount:      25,
					},
				},
				{
					Name:                  "policy2",
					FreeCapacity:          200,
					UsedCapacity:          150,
					HealthPercentage:      89.9,
					RebalancePercentage:   90.0,
					PercentageDrivesAdded: "74.3",
					RingBalance: vpsaobjectstorage.RingBalance{
						NormalPercentage:   100.0,
						DegradedPercentage: 0.0,
						CriticalPercentage: 0.0,
						NormalCount:        200,
						DegradedCount:      0,
						CriticalCount:      0,
					},
				},
			},
		},
	}, nil)

	observer.ExpectedCalls = []*mock.Call{
		// Store Metrics.
		{
			Method:    "ObserveInt64",
			Arguments: mock.Arguments{storageMetrics.AccountsCount, int64(2), mock.Anything},
		},
		{
			Method:    "ObserveInt64",
			Arguments: mock.Arguments{storageMetrics.UsersCount, int64(45), mock.Anything},
		},
		{
			Method:    "ObserveInt64",
			Arguments: mock.Arguments{storageMetrics.ContainersCount, int64(23), mock.Anything},
		},
		{
			Method:    "ObserveInt64",
			Arguments: mock.Arguments{storageMetrics.ObjectsCount, int64(78), mock.Anything},
		},
		{
			Method:    "ObserveInt64",
			Arguments: mock.Arguments{storageMetrics.DrivesCount, int64(8), mock.Anything},
		},
		{
			Method:    "ObserveInt64",
			Arguments: mock.Arguments{storageMetrics.Cache, int64(1670), mock.Anything},
		},
		// Policy 1 Metrics.
		{
			Method:    "ObserveFloat64",
			Arguments: mock.Arguments{storageMetrics.PercentageDrivesAdded, 55.2, mock.Anything},
		},
		{
			Method:    "ObserveFloat64",
			Arguments: mock.Arguments{storageMetrics.RingBalanceNormalPercentage, 75.0, mock.Anything},
		},
		{
			Method:    "ObserveFloat64",
			Arguments: mock.Arguments{storageMetrics.RingBalanceDegradedPercentage, 12.5, mock.Anything},
		},
		{
			Method:    "ObserveFloat64",
			Arguments: mock.Arguments{storageMetrics.RingBalanceCriticalPercentage, 12.5, mock.Anything},
		},
		{
			Method:    "ObserveInt64",
			Arguments: mock.Arguments{storageMetrics.FreeStorage, int64(100), mock.Anything},
		},
		{
			Method:    "ObserveInt64",
			Arguments: mock.Arguments{storageMetrics.UsedStorage, int64(50), mock.Anything},
		},
		{
			Method:    "ObserveFloat64",
			Arguments: mock.Arguments{storageMetrics.HealthPercentage, 99.9, mock.Anything},
		},
		{
			Method:    "ObserveFloat64",
			Arguments: mock.Arguments{storageMetrics.RebalancePercentage, 50.0, mock.Anything},
		},
		{
			Method:    "ObserveInt64",
			Arguments: mock.Arguments{storageMetrics.RingBalanceNormalCount, int64(150), mock.Anything},
		},
		{
			Method:    "ObserveInt64",
			Arguments: mock.Arguments{storageMetrics.RingBalanceDegradedCount, int64(25), mock.Anything},
		},
		{
			Method:    "ObserveInt64",
			Arguments: mock.Arguments{storageMetrics.RingBalanceCriticalCount, int64(25), mock.Anything},
		},
		// Policy 2 Metrics.
		{
			Method:    "ObserveFloat64",
			Arguments: mock.Arguments{storageMetrics.PercentageDrivesAdded, 74.3, mock.Anything},
		},
		{
			Method:    "ObserveFloat64",
			Arguments: mock.Arguments{storageMetrics.RingBalanceNormalPercentage, 100.0, mock.Anything},
		},
		{
			Method:    "ObserveFloat64",
			Arguments: mock.Arguments{storageMetrics.RingBalanceDegradedPercentage, 0.0, mock.Anything},
		},
		{
			Method:    "ObserveFloat64",
			Arguments: mock.Arguments{storageMetrics.RingBalanceCriticalPercentage, 0.0, mock.Anything},
		},
		{
			Method:    "ObserveInt64",
			Arguments: mock.Arguments{storageMetrics.FreeStorage, int64(200), mock.Anything},
		},
		{
			Method:    "ObserveInt64",
			Arguments: mock.Arguments{storageMetrics.UsedStorage, int64(150), mock.Anything},
		},
		{
			Method:    "ObserveFloat64",
			Arguments: mock.Arguments{storageMetrics.HealthPercentage, 89.9, mock.Anything},
		},
		{
			Method:    "ObserveFloat64",
			Arguments: mock.Arguments{storageMetrics.RebalancePercentage, 90.0, mock.Anything},
		},
		{
			Method:    "ObserveInt64",
			Arguments: mock.Arguments{storageMetrics.RingBalanceNormalCount, int64(200), mock.Anything},
		},
		{
			Method:    "ObserveInt64",
			Arguments: mock.Arguments{storageMetrics.RingBalanceDegradedCount, int64(0), mock.Anything},
		},
		{
			Method:    "ObserveInt64",
			Arguments: mock.Arguments{storageMetrics.RingBalanceCriticalCount, int64(0), mock.Anything},
		},
	}

	// Call the function being tested.
	err = storageMetrics.StorageMetricsObserve([]*config.Target{
		{
			CloudName: "cloudName",
		},
	}, func(_ context.Context, _ *config.Target) metrics.ZadaraClient {
		return mockClient
	})(context.Background(), observer)

	// Assert that there were no errors.
	require.NoError(t, err)

	// Assert that the mock client's method was called with the expected arguments.
	mockClient.AssertCalled(t, "GetAllStoragePolicies", mock.Anything, "cloudName")
}
