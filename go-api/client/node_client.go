// Package client implements the HTTP client the Go API uses to call the
// Node.js statistics API, per the required architecture: "API en Go ...
// enviará los datos resultantes a la segunda API en Node.js".
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/interseguro/matrices-go-api/matrix"
)

// NodeClient calls the Node.js statistics API.
type NodeClient struct {
	BaseURL    string
	SharedJWT  string // service-to-service token, signed with the shared secret
	HTTPClient *http.Client
}

// NewNodeClient builds a NodeClient with a sane default timeout.
func NewNodeClient(baseURL, sharedJWT string) *NodeClient {
	return &NodeClient{
		BaseURL:   baseURL,
		SharedJWT: sharedJWT,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// statsRequest mirrors the Node API's expected request body: a named set of
// matrices to aggregate statistics over.
type statsRequest struct {
	Matrices map[string]matrix.Matrix `json:"matrices"`
}

// StatsResponse mirrors the Node API's response body.
type StatsResponse struct {
	Max      float64         `json:"max"`
	Min      float64         `json:"min"`
	Average  float64         `json:"average"`
	Sum      float64         `json:"sum"`
	Diagonal map[string]bool `json:"diagonal"`
}

// FetchStats sends the named matrices to the Node.js API and returns the
// computed statistics.
func (c *NodeClient) FetchStats(ctx context.Context, matrices map[string]matrix.Matrix) (*StatsResponse, error) {
	body, err := json.Marshal(statsRequest{Matrices: matrices})
	if err != nil {
		return nil, fmt.Errorf("encoding request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/api/stats", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.SharedJWT)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling node stats API: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading node stats response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("node stats API returned %d: %s", resp.StatusCode, string(respBody))
	}

	var stats StatsResponse
	if err := json.Unmarshal(respBody, &stats); err != nil {
		return nil, fmt.Errorf("decoding node stats response: %w", err)
	}

	return &stats, nil
}
