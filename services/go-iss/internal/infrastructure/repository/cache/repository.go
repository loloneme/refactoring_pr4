package cache

import (
	"go-iss/internal/infrastructure/repository"

	"github.com/jmoiron/sqlx"
)

const (
	tableName = "space_cache"
	alias     = "sc"
	idField   = "id"
)

type Repository struct {
	*repository.Repository[int64, SpaceCache, SpaceCache]
}

func NewRepository(db *sqlx.DB) *Repository {
	cols := repository.NewColumns(readableColumns, writableColumns, alias, idField)

	baseRepo := repository.NewRepository[int64, SpaceCache, SpaceCache](db, cols, tableName)

	return &Repository{baseRepo}
}
