package models

import (
	"time"
)

type ISSFetchLog struct {
	ID        int64       `json:"id"`
	FetchedAt time.Time   `json:"fetched_at"`
	SourceURL string      `json:"source_url"`
	Payload   interface{} `json:"payload"`
}

type Trend struct {
	Movement    bool       `json:"movement"`
	DeltaKm     float64    `json:"delta_km"`
	DtSec       float64    `json:"dt_sec"`
	VelocityKmh *float64   `json:"velocity_kmh"`
	FromTime    *time.Time `json:"from_time"`
	ToTime      *time.Time `json:"to_time"`
	FromLat     *float64   `json:"from_lat"`
	FromLon     *float64   `json:"from_lon"`
	ToLat       *float64   `json:"to_lat"`
	ToLon       *float64   `json:"to_lon"`
}
