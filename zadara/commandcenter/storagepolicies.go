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

// GetAllStoragePolicies retrieves all storage policies for the client's cloud.
func (c *Client) GetAllStoragePolicies(
	ctx context.Context,
) ([]*StoreStoragePolicies, error) {
	storeRes, err := c.GetStores(ctx, c.CloudName)
	if err != nil {
		return nil, fmt.Errorf("error getting stores: %w", err)
	}

	stores := make([]*StoreStoragePolicies, len(storeRes.Zioses))

	for index, store := range storeRes.Zioses {
		stores[index] = &StoreStoragePolicies{
			Store: store,
		}

		policyRes, err := c.GetStoragePolicies(ctx, c.CloudName, store.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting storage policies: %w", err)
		}

		stores[index].Policies = policyRes.ZiosStoragePolicies
	}

	return stores, nil
}
