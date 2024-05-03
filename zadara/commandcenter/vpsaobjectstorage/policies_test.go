package vpsaobjectstorage_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/krystal/zadara-exporter/zadara/commandcenter/vpsaobjectstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetStoragePolicies(t *testing.T) {
	t.Parallel()

	// Create a mock HTTP server.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request URL.
		expectedURL := "/api/clouds/cloudName/zioses/123/storage_policies.json"
		if r.URL.Path != expectedURL {
			t.Errorf("unexpected request URL, got %s, want %s", r.URL.Path, expectedURL)
		}

		// Send a mock response.
		response := vpsaobjectstorage.ZiosStoragePoliciesResponse{
			Status: "success",
			ZiosStoragePolicies: []*vpsaobjectstorage.ZiosStoragePolicy{
				{},
				{},
			},
			Count: 2,
		}
		require.NoError(t, json.NewEncoder(w).Encode(response))
	}))
	defer server.Close()

	// Create a new client with the mock server URL.
	client := &vpsaobjectstorage.Client{
		C:       server.Client(),
		BaseURL: server.URL,
	}

	// Call the method being tested.
	ctx := context.Background()
	cloudName := "cloudName"
	ziosID := 123
	resp, err := client.GetStoragePolicies(ctx, cloudName, ziosID)
	require.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Len(t, resp.ZiosStoragePolicies, 2)
	assert.Equal(t, 2, resp.Count)
}
