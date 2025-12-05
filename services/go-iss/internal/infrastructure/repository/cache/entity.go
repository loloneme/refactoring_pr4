package cache

import (
	"time"

	"github.com/jmoiron/sqlx/types"
)

type SpaceCache struct {
	ID        int64          `db:"id"`
	Source    string         `db:"source"`
	FetchedAt time.Time      `db:"fetched_at"`
	Payload   types.JSONText `db:"payload"`
}

func (s SpaceCache) Values() []any {
	return []any{
		s.Source,
		s.FetchedAt,
		s.Payload,
	}
}

func (s SpaceCache) GetID() int64 {
	return s.ID
}
