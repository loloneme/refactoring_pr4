package fetch_and_store_iss

import (
	"context"
	"go-iss/internal/infrastructure/repository"
	"go-iss/internal/infrastructure/repository/iss"
)

type issRepo interface {
	Find(ctx context.Context, spec repository.FindSpecification) ([]iss.ISSFetchLog, error)
	Save(ctx context.Context, entity iss.ISSFetchLog) error
	GetIDFieldName(ctx context.Context) string
}

type issClient interface {
	FetchISS(ctx context.Context) (interface{}, error, int)
	GetSourceURL() string
}
