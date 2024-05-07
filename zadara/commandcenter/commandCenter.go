// Package commandcenter provides the client for the Zadara Command Centre API.
package commandcenter

import (
	"context"
	"net/http"

	"github.com/krystal/zadara-exporter/config"
	"github.com/krystal/zadara-exporter/zadara/commandcenter/vpsaobjectstorage"
)

type (
	// VPSAObjectStorage represents the VPSA Object Storage API client.
	VPSAObjectStorage interface {
		// GetStores retrieves the list of stores for the given cloudName.
		GetStores(ctx context.Context, cloudName string) (*vpsaobjectstorage.ZiosResponse, error)

		// GetStoragePolicies retrieves the storage policies for the given cloudName and ziosID.
		GetStoragePolicies(
			ctx context.Context,
			cloudName string,
			ziosID int,
		) (*vpsaobjectstorage.ZiosStoragePoliciesResponse, error)
	}

	// Client represents the client for the Zadara Command Centre API.
	Client struct {
		BaseURL   string
		CloudName string
		C         *http.Client
		VPSAObjectStorage
	}
)

// NewClient creates a new instance of the Client struct.
// It takes a pointer to a config.Target struct as a parameter and returns a pointer to the Client struct.
// The Client struct contains the necessary information to interact with the Zadara Command Centre API.
func NewClient(target *config.Target) *Client {
	httpClient := &http.Client{
		Transport: newAddTokenHeaderTransport(http.DefaultTransport, target.Token),
	}

	return &Client{
		BaseURL:           target.APIBaseURL,
		C:                 httpClient,
		CloudName:         target.CloudName,
		VPSAObjectStorage: vpsaobjectstorage.NewClient(target.APIBaseURL, httpClient),
	}
}
