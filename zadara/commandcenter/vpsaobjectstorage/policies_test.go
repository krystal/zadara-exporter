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

//nolint:funlen // most of this length is due to the test data
func TestZiosStoragePoliciesResponse(t *testing.T) {
	t.Parallel()

	testJSON := `{
		"status": "success",
		"zios_storage_policies": [
		  {
			"id": "26",
			"internal_name": "storage-policy-00000002",
			"name": "3-way-protection-policy",
			"status": "initialised",
			"protection": "3-way",
			"default": null,
			"health_status": "normal",
			"health_percentage": 100,
			"ring_balance": {
			  "normal_percentage": 99,
			  "degraded_percentage": 0,
			  "critical_percentage": 0,
			  "normal_count": 32729,
			  "degraded_count": 39,
			  "critical_count": 0
			},
			"used_capacity": 0,
			"free_capacity": 299573968896
		  },
		  {
			"name": "MetadataPolicy",
			"status": "initialised",
			"protection": "3-way",
			"default": null,
			"health_status": "normal",
			"health_percentage": 100,
			"ring_balance": {
			  "normal_percentage": 84,
			  "degraded_percentage": 15,
			  "critical_percentage": 0,
			  "normal_count": 3465,
			  "degraded_count": 631,
			  "critical_count": 0
			},
			"used_capacity": 38377881,
			"free_capacity": 236184823399
		  }
		],
		"count": 2
	  }`

	var resp vpsaobjectstorage.ZiosStoragePoliciesResponse

	require.NoError(t, json.Unmarshal([]byte(testJSON), &resp))

	assert.Equal(t, "success", resp.Status)

	assert.Len(t, resp.ZiosStoragePolicies, 2)

	assert.Equal(t, &vpsaobjectstorage.ZiosStoragePolicy{
		ID:               "26",
		InternalName:     "storage-policy-00000002",
		Name:             "3-way-protection-policy",
		Status:           "initialised",
		Protection:       "3-way",
		Default:          nil,
		HealthStatus:     "normal",
		HealthPercentage: 100,
		RingBalance: vpsaobjectstorage.RingBalance{
			NormalPercentage:   99,
			DegradedPercentage: 0,
			CriticalPercentage: 0,
			NormalCount:        32729,
			DegradedCount:      39,
			CriticalCount:      0,
		},
		UsedCapacity: 0,
		FreeCapacity: 299573968896,
	}, resp.ZiosStoragePolicies[0])

	assert.Equal(t, &vpsaobjectstorage.ZiosStoragePolicy{
		Name:             "MetadataPolicy",
		Status:           "initialised",
		Protection:       "3-way",
		Default:          nil,
		HealthStatus:     "normal",
		HealthPercentage: 100,
		RingBalance: vpsaobjectstorage.RingBalance{
			NormalPercentage:   84,
			DegradedPercentage: 15,
			CriticalPercentage: 0,
			NormalCount:        3465,
			DegradedCount:      631,
			CriticalCount:      0,
		},
		UsedCapacity: 38377881,
		FreeCapacity: 236184823399,
	}, resp.ZiosStoragePolicies[1])

	assert.Equal(t, 2, resp.Count)
}
