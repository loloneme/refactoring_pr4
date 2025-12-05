package nasa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	APIKey         string `env:"NASA_API_KEY" envDefault:""`
	TimeoutSeconds int    `env:"NASA_CLIENT_TIMEOUT_SECONDS" envDefault:"30"`
	OSDRURL        string `env:"NASA_OSDR_URL" envDefault:"https://visualization.osdr.nasa.gov/biodata/api/v2/datasets/?format=json"`
	APODURL        string `env:"NASA_APOD_URL" envDefault:"https://api.nasa.gov/planetary/apod"`
	NEOFeedURL     string `env:"NASA_NEO_FEED_URL" envDefault:"https://api.nasa.gov/neo/rest/v1/feed"`
	DONKIFLRURL    string `env:"NASA_DONKI_FLR_URL" envDefault:"https://api.nasa.gov/DONKI/FLR"`
	DONKICMEURL    string `env:"NASA_DONKI_CME_URL" envDefault:"https://api.nasa.gov/DONKI/CME"`
}

type Client struct {
	apiKey      string
	timeout     time.Duration
	client      *http.Client
	osdrURL     string
	apodURL     string
	neoFeedURL  string
	donkiFLRURL string
	donkiCMEURL string
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
		return nil, fmt.Errorf("failed to load NASA Client config: %w", err)
	}

	timeout := time.Duration(cfg.TimeoutSeconds) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &Client{
		apiKey:      cfg.APIKey,
		timeout:     timeout,
		client:      &http.Client{Timeout: timeout},
		osdrURL:     cfg.OSDRURL,
		apodURL:     cfg.APODURL,
		neoFeedURL:  cfg.NEOFeedURL,
		donkiFLRURL: cfg.DONKIFLRURL,
		donkiCMEURL: cfg.DONKICMEURL,
	}, nil
}

func (c *Client) FetchOSDR(ctx context.Context, osdrURL string) (interface{}, error) {
	if osdrURL == "" {
		osdrURL = c.osdrURL
	}
	req, err := http.NewRequestWithContext(ctx, "GET", osdrURL, nil)
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

func (c *Client) GetOSDRURL() string {
	return c.osdrURL
}

func (c *Client) Close() {
	c.client.CloseIdleConnections()
}

func (c *Client) FetchAPOD(ctx context.Context) (interface{}, error) {
	apiURL := c.apodURL
	reqURL, err := url.Parse(apiURL)
	if err != nil {
		return nil, fmt.Errorf("parse URL: %w", err)
	}

	q := reqURL.Query()
	q.Set("thumbs", "true")
	if c.apiKey != "" {
		q.Set("api_key", c.apiKey)
	}
	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	var jsonData interface{}
	if err := json.NewDecoder(resp.Body).Decode(&jsonData); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return jsonData, nil
}

// FetchNEOFeed получает данные из NEO Feed API
func (c *Client) FetchNEOFeed(ctx context.Context, startDate, endDate string) (interface{}, error) {
	apiURL := c.neoFeedURL
	reqURL, err := url.Parse(apiURL)
	if err != nil {
		return nil, fmt.Errorf("parse URL: %w", err)
	}

	q := reqURL.Query()
	q.Set("start_date", startDate)
	q.Set("end_date", endDate)
	if c.apiKey != "" {
		q.Set("api_key", c.apiKey)
	}
	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	var jsonData interface{}
	if err := json.NewDecoder(resp.Body).Decode(&jsonData); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return jsonData, nil
}

// FetchDONKIFLR получает данные из DONKI FLR API
func (c *Client) FetchDONKIFLR(ctx context.Context, startDate, endDate string) (interface{}, error) {
	apiURL := c.donkiFLRURL
	reqURL, err := url.Parse(apiURL)
	if err != nil {
		return nil, fmt.Errorf("parse URL: %w", err)
	}

	q := reqURL.Query()
	q.Set("startDate", startDate)
	q.Set("endDate", endDate)
	if c.apiKey != "" {
		q.Set("api_key", c.apiKey)
	}
	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	var jsonData interface{}
	if err := json.NewDecoder(resp.Body).Decode(&jsonData); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return jsonData, nil
}

// FetchDONKICME получает данные из DONKI CME API
func (c *Client) FetchDONKICME(ctx context.Context, startDate, endDate string) (interface{}, error) {
	apiURL := c.donkiCMEURL
	reqURL, err := url.Parse(apiURL)
	if err != nil {
		return nil, fmt.Errorf("parse URL: %w", err)
	}

	q := reqURL.Query()
	q.Set("startDate", startDate)
	q.Set("endDate", endDate)
	if c.apiKey != "" {
		q.Set("api_key", c.apiKey)
	}
	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	var jsonData interface{}
	if err := json.NewDecoder(resp.Body).Decode(&jsonData); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return jsonData, nil
}
