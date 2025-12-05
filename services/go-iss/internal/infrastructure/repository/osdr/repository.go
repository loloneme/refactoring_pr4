package osdr

import (
	"go-iss/internal/infrastructure/repository"

	"github.com/jmoiron/sqlx"
)

const (
	tableName = "osdr_items"
	alias     = "oi"
	idField   = "id"
)

type Repository struct {
	*repository.Repository[int64, OSDRItem, OSDRItem]
}

func NewRepository(db *sqlx.DB) *Repository {
	cols := repository.NewColumns(readableColumns, writableColumns, alias, idField)

	baseRepo := repository.NewRepository[int64, OSDRItem, OSDRItem](db, cols, tableName)

	return &Repository{
		Repository: baseRepo,
	}
}
