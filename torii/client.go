package torii

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

const baseURL = "https://api.toriihq.com"

// Client wraps the Torii REST API.
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// newClient creates a new Client using the given API key.
func newClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// getClient retrieves a configured Client from the plugin connection config.
func getClient(ctx context.Context, d *plugin.QueryData) (*Client, error) {
	cfg := GetConfig(d.Connection)
	if cfg.APIKey == nil || *cfg.APIKey == "" {
		return nil, fmt.Errorf("api_key must be configured")
	}
	return newClient(*cfg.APIKey), nil
}

// get performs an authenticated GET request and unmarshals the response body into result.
func (c *Client) get(ctx context.Context, path string, params map[string]string, result interface{}) error {
	u, err := url.Parse(baseURL + path)
	if err != nil {
		return fmt.Errorf("parsing URL: %w", err)
	}

	if len(params) > 0 {
		q := u.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("unmarshalling response: %w", err)
	}

	return nil
}
