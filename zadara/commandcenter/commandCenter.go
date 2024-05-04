// Package commandcenter provides the client for the Zadara Command Centre API.
package commandcenter

import (
	"context"
	"net/http"

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
		BaseURL string
		C       *http.Client
		VPSAObjectStorage
	}
)

// NewClientFromToken creates a new client for interacting with the Zadara Command
// Centre API using the provided API token.
// It takes the base URL of the Command Centre API and the API token as parameters.
// It returns a pointer to the created Client.
func NewClientFromToken(baseURL, apiToken string) *Client {
	httpClient := &http.Client{
		Transport: newAddTokenHeaderTransport(http.DefaultTransport, apiToken),
	}

	return &Client{
		BaseURL:           baseURL,
		C:                 httpClient,
		VPSAObjectStorage: vpsaobjectstorage.NewClient(baseURL, httpClient),
	}
}

// NewClient creates a new instance of the Client struct.
// It initialises the Client with the provided baseURL and http.Client.
// It also initialises the VPSAObjectStorage field with a new instance of the vpsaobjectstorage.Client.
func NewClient(baseURL string, c *http.Client) *Client {
	return &Client{
		BaseURL:           baseURL,
		C:                 c,
		VPSAObjectStorage: vpsaobjectstorage.NewClient(baseURL, c),
	}
}
