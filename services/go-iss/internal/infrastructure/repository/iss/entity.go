package iss

import (
	"time"

	"github.com/jmoiron/sqlx/types"
)

type ISSFetchLog struct {
	ID        int64          `db:"id"`
	FetchedAt time.Time      `db:"fetched_at"`
	SourceURL string         `db:"source_url"`
	Payload   types.JSONText `db:"payload"`
}

func (i ISSFetchLog) Values() []any {
	return []any{
		i.FetchedAt,
		i.SourceURL,
		i.Payload,
	}
}

func (i ISSFetchLog) GetID() int64 {
	return i.ID
}
