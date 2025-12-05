package models

import "time"

type SpaceRefreshResponse struct {
	Refreshed []string `json:"refreshed"`
}

type SpaceLatestResponse struct {
	Source    string      `json:"source,omitempty"`
	FetchedAt time.Time   `json:"fetched_at"`
	Payload   interface{} `json:"payload"`
}

type SpaceLatestNoDataResponse struct {
	Source  string `json:"source"`
	Message string `json:"message"`
}

type SpaceSummary struct {
	APOD      SpaceCacheItem `json:"apod"`
	NEO       SpaceCacheItem `json:"neo"`
	FLR       SpaceCacheItem `json:"flr"`
	CME       SpaceCacheItem `json:"cme"`
	SpaceX    SpaceCacheItem `json:"spacex"`
	ISS       SpaceCacheItem `json:"iss"`
	OSDRCount int64          `json:"osdr_count"`
}

type SpaceCacheItem struct {
	At      *time.Time  `json:"at,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}
