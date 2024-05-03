package vpsaobjectstorage

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"strconv"
)

type (
	// RingBalance represents the balance of the ring.
	RingBalance struct {
		NormalPercentage   int `json:"normal_percentage"`
		DegradedPercentage int `json:"degraded_percentage"`
		CriticalPercentage int `json:"critical_percentage"`
		NormalCount        int `json:"normal_count"`
		DegradedCount      int `json:"degraded_count"`
		CriticalCount      int `json:"critical_count"`
	}

	// ZiosStoragePolicy represents a VPSA Object Storage storage policy.
	ZiosStoragePolicy struct {
		ID               string      `json:"id"`
		InternalName     string      `json:"internal_name"`
		Name             string      `json:"name"`
		Status           string      `json:"status"`
		Protection       string      `json:"protection"`
		Default          interface{} `json:"default"`
		HealthStatus     string      `json:"health_status"`
		HealthPercentage int         `json:"health_percentage"`
		RingBalance      RingBalance `json:"ring_balance"`
		UsedCapacity     int64       `json:"used_capacity"`
		FreeCapacity     int64       `json:"free_capacity"`
	}

	// ZiosStoragePoliciesResponse represents the response of the GetStoragePolicies API.
	ZiosStoragePoliciesResponse struct {
		Status              string               `json:"status"`
		ZiosStoragePolicies []*ZiosStoragePolicy `json:"zios_storage_policies"`
		Count               int                  `json:"count"`
	}
)

// GetStoragePolicies returns the list of a VPSA Object Storage storage policies.
//
// API Docs
// Returns the list of a VPSA Object Storage storage policies.
// GET /api/clouds/{cloud_name}/zioses/{id or internal-name}/storage_policies(.xml/json)
// Example:
// curl -X GET -H "Content-Type: application/json" -H "X-Token: <token>" \
// 'https://<command-center-ip>:8888/api/clouds/{cloud_name}/zioses/{id or internal-name}/storage_policies.json'.
func (c *Client) GetStoragePolicies(
	ctx context.Context,
	cloudName string, ziosID int,
) (*ZiosStoragePoliciesResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		c.BaseURL+path.Join("/api/clouds", cloudName, "zioses", strconv.Itoa(ziosID), "storage_policies.json"),
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

	var resp ZiosStoragePoliciesResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &resp, nil
}
