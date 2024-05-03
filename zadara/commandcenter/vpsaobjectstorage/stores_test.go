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

func TestClient_GetStores(t *testing.T) {
	t.Parallel()

	// Create a mock HTTP server.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request URL.
		expectedURL := "/api/clouds/cloudName/zioses.json"
		if r.URL.Path != expectedURL {
			t.Errorf("unexpected request URL, got %s, want %s", r.URL.Path, expectedURL)
		}

		// Send a mock response.
		response := vpsaobjectstorage.ZiosResponse{
			Status: "success",
			Zioses: []*vpsaobjectstorage.Zios{
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
	resp, err := client.GetStores(ctx, cloudName)
	require.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Len(t, resp.Zioses, 2)
	assert.Equal(t, 2, resp.Count)
}
