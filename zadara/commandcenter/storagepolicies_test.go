package commandcenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/krystal/zadara-exporter/zadara/commandcenter"
	"github.com/krystal/zadara-exporter/zadara/commandcenter/vpsaobjectstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type (
	// MockClient is a mock implementation of the vpsaobjectstorage.Client interface.
	MockClient struct {
		mock.Mock
	}
)

func (m *MockClient) GetStores(ctx context.Context, cloudName string) (*vpsaobjectstorage.ZiosResponse, error) {
	args := m.Called(ctx, cloudName)

	firstarg, ok := args.Get(0).(*vpsaobjectstorage.ZiosResponse)

	var err error
	if args.Error(1) != nil {
		err = fmt.Errorf("error returned from getStores: %w", args.Error(1))
	}

	if !ok {
		return nil, err
	}

	return firstarg, err
}

func (m *MockClient) GetStoragePolicies(
	ctx context.Context,
	cloudName string,
	ziosID int,
) (*vpsaobjectstorage.ZiosStoragePoliciesResponse, error) {
	args := m.Called(ctx, cloudName, ziosID)

	firstarg, ok := args.Get(0).(*vpsaobjectstorage.ZiosStoragePoliciesResponse)

	var err error
	if args.Error(1) != nil {
		err = fmt.Errorf("error returned from getStores: %w", args.Error(1))
	}

	if !ok {
		return nil, err
	}

	return firstarg, err
}

func TestClient_GetAllStoragePolicies(t *testing.T) {
	t.Parallel()

	// Create a mock client.
	mockClient := new(MockClient)

	// Set up the expected calls and responses.
	cloudName := "cloudName"
	storeRes := &vpsaobjectstorage.ZiosResponse{
		Zioses: []*vpsaobjectstorage.Zios{
			{ID: 1},
			{ID: 2},
		},
	}
	policyRes := &vpsaobjectstorage.ZiosStoragePoliciesResponse{
		ZiosStoragePolicies: []*vpsaobjectstorage.ZiosStoragePolicy{
			{ID: 1},
			{ID: 2},
		},
	}

	mockClient.On("GetStores", mock.Anything, cloudName).Return(storeRes, nil).Once()
	mockClient.On("GetStoragePolicies", mock.Anything, cloudName, 1).Return(policyRes, nil).Once()
	mockClient.On("GetStoragePolicies", mock.Anything, cloudName, 2).Return(policyRes, nil).Once()

	// Create the client under test.
	client := commandcenter.Client{
		CloudName:         cloudName,
		VPSAObjectStorage: mockClient,
	}

	// Call the method being tested.
	ctx := context.Background()
	stores, err := client.GetAllStoragePolicies(ctx)

	// Assert the results.
	require.NoError(t, err)
	assert.Len(t, stores, 2)
	assert.Equal(t, 1, stores[0].Policies[0].ID)
	assert.Equal(t, 2, stores[0].Policies[1].ID)
	assert.Equal(t, 1, stores[1].Policies[0].ID)
	assert.Equal(t, 2, stores[1].Policies[1].ID)

	// Verify the expected calls were made.
	mockClient.AssertExpectations(t)
}
