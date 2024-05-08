// Package vpsaobjectstorage provides the client for the VPSA Object Storage API.
package vpsaobjectstorage

import (
	"net/http"
)

type (
	// Client represents the client for the VPSA Object Storage API.
	Client struct {
		BaseURL   string
		C         *http.Client
		CloudName string
	}
)

// NewClient returns a new VPSA Object Storage API client.
func NewClient(baseURL string, c *http.Client) *Client {
	return &Client{
		BaseURL: baseURL,
		C:       c,
	}
}
