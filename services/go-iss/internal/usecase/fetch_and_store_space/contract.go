package fetch_and_store_space

import (
	"context"
	"go-iss/internal/infrastructure/repository/cache"
	"go-iss/internal/rpc/errors"
)

type RefreshSpaceService interface {
	RefreshSpace(ctx context.Context, sources []string) ([]string, *errors.ServiceError)
}

type cacheRepo interface {
	Save(ctx context.Context, entity cache.SpaceCache) error
}

type nasaClient interface {
	FetchAPOD(ctx context.Context) (interface{}, error)
	FetchNEOFeed(ctx context.Context, startDate, endDate string) (interface{}, error)
	FetchDONKIFLR(ctx context.Context, startDate, endDate string) (interface{}, error)
	FetchDONKICME(ctx context.Context, startDate, endDate string) (interface{}, error)
}

type spacexClient interface {
	FetchNextLaunch(ctx context.Context) (interface{}, error)
}
