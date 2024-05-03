package vpsaobjectstorage

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"path"
)

type (
	// NetworkConfiguration represents the network configuration of the Zios.
	NetworkConfiguration struct {
		FeIP string `json:"fe_ip"`
		BeIP string `json:"be_ip"`
		HbIP string `json:"hb_ip"`
	}

	// Zios represents a VPSA Object Storage object store.
	Zios struct {
		ID                    int                             `json:"id"`
		Name                  string                          `json:"name"`
		InternalName          string                          `json:"internal_name"`
		User                  string                          `json:"user"`
		TenantName            string                          `json:"tenant_name"`
		Description           string                          `json:"description"`
		Status                string                          `json:"status"`
		EngineType            string                          `json:"engine_type"`
		Vcpus                 int                             `json:"vcpus"`
		RAM                   int                             `json:"ram"`
		HTTPSTermination      bool                            `json:"https_termination"`
		Image                 string                          `json:"image"`
		Drives                int                             `json:"drives"`
		Cache                 int                             `json:"cache"`
		VirtualControllers    int                             `json:"virtual_controllers"`
		IPAddress             string                          `json:"ip_address"`
		PublicIP              interface{}                     `json:"public_ip"`
		ManagementURL         string                          `json:"management_url"`
		StoragePoliciesCount  int                             `json:"storage_policies_count"`
		MetadataPoliciesCount int                             `json:"metadata_policies_count"`
		AccountsCount         int                             `json:"accounts_count"`
		UsersCount            int                             `json:"users_count"`
		ContainersCount       int                             `json:"containers_count"`
		ObjectsCount          int                             `json:"objects_count"`
		NetworkConfiguration  map[string]NetworkConfiguration `json:"network_configuration"`
		CreatedAt             string                          `json:"created_at"`
		UpdatedAt             string                          `json:"updated_at"`
	}

	// ZiosResponse represents the response of the GetStores API.
	ZiosResponse struct {
		Status string  `json:"status"`
		Zioses []*Zios `json:"zioses"`
		Count  int     `json:"count"`
	}
)

// GetStores returns a list of all VPSA Object Storage object stores.
//
// API Docs
// Returns a list of all VPSA Object Storage object stores.
// GET /api/clouds/{cloud_name}/zioses(.xml/json)
// Example:
// curl -X GET -H "Content-Type: application/json" -H "X-Token: <token>" \
// 'https://<command-center-ip>:8888/api/clouds/{cloud_name}/zioses.json?page=1&per_page=10'
// page	Integer	The page number to start from.
// per_page	Integer	The total number of records to return.
func (c *Client) GetStores(
	ctx context.Context,
	cloudName string,
) (*ZiosResponse, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		c.BaseURL+path.Join("/api/clouds", cloudName, "zioses.json"),
		nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	res, err := c.C.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error in sending request: %w", err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			slog.Error("error closing response body", "error", err)
		}
	}()

	var resp ZiosResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &resp, nil
}