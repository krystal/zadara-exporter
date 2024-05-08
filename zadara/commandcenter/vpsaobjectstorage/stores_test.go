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

//nolint:funlen // most of this length is due to the test data
func TestZiosResponse(t *testing.T) {
	t.Parallel()

	testJSON := `{
		"status": "success",
		"zioses": [
		  {
			"id": 188,
			"name": "zios",
			"internal_name": "vsa-0000016d",
			"user": "john",
			"tenant_name": "user_john_Gt1cj",
			"description": "zios",
			"status": "Normal",
			"engine_type": "ZIOS",
			"vcpus": 3,
			"ram": 6144,
			"https_termination": true,
			"image": "zios-00.00-434.img",
			"drives": 5,
			"cache": 30,
			"virtual_controllers": 3,
			"ip_address": "150.50.2.130",
			"public_ip": null,
			"management_url": "vsa-0000016d-zadara-dev2.zadarazios.com",
			"storage_policies_count": 1,
			"metadata_policies_count": 1,
			"accounts_count": 1,
			"users_count": 3,
			"containers_count": 0,
			"objects_count": 0,
			"network_configuration": {
			  "vc0": {
				"fe_ip": "150.50.2.113",
				"be_ip": "150.50.3.113",
				"hb_ip": "150.50.0.113"
			  },
			  "vc1": {
				"fe_ip": "150.50.2.114",
				"be_ip": "150.50.3.114",
				"hb_ip": "150.50.0.114"
			  },
			  "vc2": {
				"fe_ip": "150.50.2.118",
				"be_ip": "150.50.3.118",
				"hb_ip": "150.50.0.118"
			  }
			},
			"created_at": "2016-04-15 20:22:10 UTC",
			"updated_at": "2016-04-18 13:42:48 UTC"
		  }
		],
		"count": 1
	  }`

	var resp vpsaobjectstorage.ZiosResponse

	require.NoError(t, json.Unmarshal([]byte(testJSON), &resp))

	assert.Equal(t, "success", resp.Status)

	assert.Len(t, resp.Zioses, 1)

	assert.Equal(t, &vpsaobjectstorage.Zios{
		ID:                    188,
		Name:                  "zios",
		InternalName:          "vsa-0000016d",
		User:                  "john",
		TenantName:            "user_john_Gt1cj",
		Description:           "zios",
		Status:                "Normal",
		EngineType:            "ZIOS",
		Vcpus:                 3,
		RAM:                   6144,
		HTTPSTermination:      true,
		Image:                 "zios-00.00-434.img",
		Drives:                5,
		Cache:                 30,
		VirtualControllers:    3,
		IPAddress:             "150.50.2.130",
		PublicIP:              nil,
		ManagementURL:         "vsa-0000016d-zadara-dev2.zadarazios.com",
		StoragePoliciesCount:  1,
		MetadataPoliciesCount: 1,
		AccountsCount:         1,
		UsersCount:            3,
		ContainersCount:       0,
		ObjectsCount:          0,
		NetworkConfiguration: map[string]vpsaobjectstorage.NetworkConfiguration{
			"vc0": {
				FeIP: "150.50.2.113",
				BeIP: "150.50.3.113",
				HbIP: "150.50.0.113",
			},
			"vc1": {
				FeIP: "150.50.2.114",
				BeIP: "150.50.3.114",
				HbIP: "150.50.0.114",
			},
			"vc2": {
				FeIP: "150.50.2.118",
				BeIP: "150.50.3.118",
				HbIP: "150.50.0.118",
			},
		},
		CreatedAt: "2016-04-15 20:22:10 UTC",
		UpdatedAt: "2016-04-18 13:42:48 UTC",
	}, resp.Zioses[0])

	assert.Equal(t, 1, resp.Count)
}
