package commandcenter

import (
	"context"
	"fmt"

	"github.com/krystal/zadara-exporter/zadara/commandcenter/vpsaobjectstorage"
)

// GetAllStoragePolicies : gets all storage policies for a cloud.
func (c *Client) GetAllStoragePolicies(
	ctx context.Context,
	cloudName string,
) ([]*vpsaobjectstorage.ZiosStoragePolicy, error) {
	storeRes, err := c.GetStores(ctx, cloudName)
	if err != nil {
		return nil, fmt.Errorf("error getting stores: %w", err)
	}

	policies := []*vpsaobjectstorage.ZiosStoragePolicy{}

	for _, store := range storeRes.Zioses {
		policyRes, err := c.GetStoragePolicies(ctx, cloudName, store.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting storage policies: %w", err)
		}

		policies = append(policies, policyRes.ZiosStoragePolicies...)
	}

	return policies, nil
}
