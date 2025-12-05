package get_space_summary

import (
	"context"
	"go-iss/internal/infrastructure/repository"
	"go-iss/internal/infrastructure/repository/cache"
	"go-iss/internal/infrastructure/repository/iss"
)

type cacheRepo interface {
	Find(ctx context.Context, spec repository.FindSpecification) ([]cache.SpaceCache, error)
	GetIDFieldName(ctx context.Context) string
}

type issRepo interface {
	Find(ctx context.Context, spec repository.FindSpecification) ([]iss.ISSFetchLog, error)
	GetIDFieldName(ctx context.Context) string
}

type osdrRepo interface {
	Count(ctx context.Context) (int64, error)
}
