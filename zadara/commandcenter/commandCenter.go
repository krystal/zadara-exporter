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
		GetStores(ctx context.Context, cloudName string) (*vpsaobjectstorage.ZiosResponse, error)
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

// NewClientFromToken returns a new Zadara Command Centre API client with the given API token.
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

// NewClient returns a new Zadara Command Centre API client using the given HTTP client.
func NewClient(baseURL string, c *http.Client) *Client {
	return &Client{
		BaseURL:           baseURL,
		C:                 c,
		VPSAObjectStorage: vpsaobjectstorage.NewClient(baseURL, c),
	}
}
