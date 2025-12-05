package iss

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	BaseURL        string `env:"WHERE_ISS_URL" envDefault:"https://api.wheretheiss.at/v1/satellites/25544"`
	TimeoutSeconds int    `env:"ISS_CLIENT_TIMEOUT_SECONDS" envDefault:"20"`
}

type Client struct {
	BaseURL string
	timeout time.Duration
	client  *http.Client
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func NewFromConfig() (*Client, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load ISS client config: %w", err)
	}

	timeout := time.Duration(cfg.TimeoutSeconds) * time.Second
	if timeout == 0 {
		timeout = 20 * time.Second
	}

	return &Client{
		BaseURL: cfg.BaseURL,
		timeout: timeout,
		client: &http.Client{
			Timeout: timeout,
		},
	}, nil
}

func (c *Client) Close() {
	c.client.CloseIdleConnections()
}

func (c *Client) GetSourceURL() string {
	return c.BaseURL
}

func (c *Client) FetchISS(ctx context.Context) (interface{}, error, int) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.BaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err), 0
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err), 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode), resp.StatusCode
	}

	var jsonData interface{}
	if err := json.NewDecoder(resp.Body).Decode(&jsonData); err != nil {
		return nil, fmt.Errorf("decode response: %w", err), resp.StatusCode
	}

	return jsonData, nil, resp.StatusCode
}
