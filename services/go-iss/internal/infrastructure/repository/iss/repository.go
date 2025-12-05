package iss

import (
	"go-iss/internal/infrastructure/repository"

	"github.com/jmoiron/sqlx"
)

const (
	tableName = "iss_fetch_log"
	alias     = "ifl"
	idField   = "id"
)

type Repository struct {
	*repository.Repository[int64, ISSFetchLog, ISSFetchLog]
}

func NewRepository(db *sqlx.DB) *Repository {
	cols := repository.NewColumns(readableColumns, writableColumns, alias, idField)

	baseRepo := repository.NewRepository[int64, ISSFetchLog, ISSFetchLog](db, cols, tableName)

	return &Repository{baseRepo}
}
