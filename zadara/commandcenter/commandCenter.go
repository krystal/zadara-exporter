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
		BaseURL   string
		CloudName string
		C         *http.Client
		VPSAObjectStorage
	}
)

// NewClientFromToken creates a new client with the provided base URL, API token, and cloud name.
// It returns a pointer to the created Client.
func NewClientFromToken(baseURL, apiToken, cloudName string) *Client {
	httpClient := &http.Client{
		Transport: newAddTokenHeaderTransport(http.DefaultTransport, apiToken),
	}

	return &Client{
		BaseURL:           baseURL,
		C:                 httpClient,
		CloudName:         cloudName,
		VPSAObjectStorage: vpsaobjectstorage.NewClient(baseURL, httpClient),
	}
}

// NewClient creates a new instance of the Client struct.
// It initialises the Client with the provided baseURL, http.Client, and cloudName.
// It also creates a new instance of the vpsaobjectstorage.Client and assigns
// it to the VPSAObjectStorage field of the Client.
func NewClient(baseURL string, c *http.Client, cloudName string) *Client {
	return &Client{
		BaseURL:           baseURL,
		C:                 c,
		CloudName:         cloudName,
		VPSAObjectStorage: vpsaobjectstorage.NewClient(baseURL, c),
	}
}
