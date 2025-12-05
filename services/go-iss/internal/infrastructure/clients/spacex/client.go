package spacex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	BaseURL        string `env:"SPACEX_API_URL" envDefault:"https://api.spacexdata.com/v4/launches/next"`
	TimeoutSeconds int    `env:"SPACEX_CLIENT_TIMEOUT_SECONDS" envDefault:"30"`
}

type Client struct {
	baseURL string
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
		return nil, fmt.Errorf("failed to load SpaceX client config: %w", err)
	}

	timeout := time.Duration(cfg.TimeoutSeconds) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &Client{
		baseURL: cfg.BaseURL,
		timeout: timeout,
		client: &http.Client{
			Timeout: timeout,
		},
	}, nil
}

func (c *Client) Close() {
	c.client.CloseIdleConnections()
}

func (c *Client) FetchNextLaunch(ctx context.Context) (interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var jsonData interface{}
	if err := json.NewDecoder(resp.Body).Decode(&jsonData); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return jsonData, nil
}
