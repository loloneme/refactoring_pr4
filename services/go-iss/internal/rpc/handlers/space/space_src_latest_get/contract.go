package space_src_latest_get

import (
	"context"
	"go-iss/internal/infrastructure/repository/cache"
	"go-iss/internal/rpc/errors"
)

type GetLatestSpaceCacheService interface {
	GetLatestSpaceCache(ctx context.Context, source string) (*cache.SpaceCache, *errors.ServiceError)
}
