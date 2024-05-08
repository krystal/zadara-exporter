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
		NormalPercentage   float64 `json:"normal_percentage"`
		DegradedPercentage float64 `json:"degraded_percentage"`
		CriticalPercentage float64 `json:"critical_percentage"`
		NormalCount        int64   `json:"normal_count"`
		DegradedCount      int64   `json:"degraded_count"`
		CriticalCount      int64   `json:"critical_count"`
	}

	// ZiosStoragePolicy represents a VPSA Object Storage storage policy.
	ZiosStoragePolicy struct {
		ID                                    int         `json:"id"`
		InternalName                          string      `json:"internal_name"`
		Name                                  string      `json:"name"`
		Status                                string      `json:"status"`
		Protection                            string      `json:"protection"`
		RebalanceCurrentCompletionProjectedAt string      `json:"rebalance_current_completion_projected_at"`
		RebalancePercentage                   float64     `json:"rebalance_percentage"`
		PercentageDrivesAdded                 string      `json:"percentage_drives_added"`
		RebalancingPaused                     bool        `json:"rebalancing_paused"`
		Default                               bool        `json:"default"`
		HealthStatus                          string      `json:"health_status"`
		HealthPercentage                      float64     `json:"health_percentage"`
		RingBalance                           RingBalance `json:"ring_balance"`
		UsedCapacity                          int64       `json:"used_capacity"`
		FreeCapacity                          int64       `json:"free_capacity"`
	}

	// ZiosStoragePoliciesResponse represents the response of the GetStoragePolicies API.
	ZiosStoragePoliciesResponse struct {
		Status              string               `json:"status"`
		Message             string               `json:"message"`
		ZiosStoragePolicies []*ZiosStoragePolicy `json:"zios_storage_policies"`
		Count               int                  `json:"count"`
	}
)

// GetStoragePolicies retrieves the storage policies for a specific Zios object in a cloud.
// It takes a context, cloud name, and Zios ID as input parameters.
// It returns a pointer to a ZiosStoragePoliciesResponse struct and an error.
// The ZiosStoragePoliciesResponse struct contains the response data from the API call.
// If there is an error creating the request, sending the request, closing the response body,
// or decoding the response, an error is returned.
//
// # API Docs
//
// Returns the list of a VPSA Object Storage storage policies.
// GET /api/clouds/{cloud_name}/zioses/{id or internal-name}/storage_policies(.xml/json)
//
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

	if res.StatusCode != http.StatusOK || resp.Status == "error" {
		return nil, fmt.Errorf("%w: %s", ErrResponse, resp.Message)
	}

	return &resp, nil
}
