package commandcenter

import (
	"context"
	"fmt"

	"github.com/krystal/zadara-exporter/zadara/commandcenter/vpsaobjectstorage"
)

type (
	// StoreStoragePolicies represents a store and its associated storage policies.
	StoreStoragePolicies struct {
		Store    *vpsaobjectstorage.Zios
		Policies []*vpsaobjectstorage.ZiosStoragePolicy
	}
)

// GetAllStoragePolicies retrieves all storage policies for a given cloud.
// It returns a slice of ZiosStoragePolicy objects and an error, if any.
func (c *Client) GetAllStoragePolicies(
	ctx context.Context,
	cloudName string,
) ([]*StoreStoragePolicies, error) {
	storeRes, err := c.GetStores(ctx, c.CloudName)
	if err != nil {
		return nil, fmt.Errorf("error getting stores: %w", err)
	}

	stores := make([]*StoreStoragePolicies, len(storeRes.Zioses))

	for index, store := range storeRes.Zioses {
		stores[index].Store = store

		policyRes, err := c.GetStoragePolicies(ctx, cloudName, store.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting storage policies: %w", err)
		}

		stores[index].Policies = policyRes.ZiosStoragePolicies
	}

	return stores, nil
}
