package osdr

import (
	"time"

	"github.com/jmoiron/sqlx/types"
)

type OSDRItem struct {
	ID         int64          `db:"id"`
	DatasetID  *string        `db:"dataset_id"`
	Title      *string        `db:"title"`
	Status     *string        `db:"status"`
	UpdatedAt  *time.Time     `db:"updated_at"`
	InsertedAt time.Time      `db:"inserted_at"`
	Raw        types.JSONText `db:"raw"`
}

func (o OSDRItem) Values() []any {
	return []any{
		o.DatasetID,
		o.Title,
		o.Status,
		o.UpdatedAt,
		o.InsertedAt,
		o.Raw,
	}
}

func (o OSDRItem) GetID() int64 {
	return o.ID
}
