package get_latest_space_cache

import (
	"context"
	"go-iss/internal/infrastructure/repository"
	"go-iss/internal/infrastructure/repository/cache"
)

type GetLatestSpaceCacheService interface {
	GetLatestSpaceCache(ctx context.Context, source string) (*cache.SpaceCache, error)
}

type cacheRepo interface {
	Find(ctx context.Context, spec repository.FindSpecification) ([]cache.SpaceCache, error)
	GetIDFieldName(ctx context.Context) string
}
