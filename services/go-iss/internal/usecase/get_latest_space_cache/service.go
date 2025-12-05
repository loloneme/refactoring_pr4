package get_latest_space_cache

import (
	"context"
	"errors"
	"go-iss/internal/domain/space_cache/specification"
	"go-iss/internal/infrastructure/repository"
	"go-iss/internal/infrastructure/repository/cache"
	rpc_errors "go-iss/internal/rpc/errors"
)

type Service struct {
	cacheRepo cacheRepo
}

func New(repo cacheRepo) *Service {
	return &Service{
		cacheRepo: repo,
	}
}

func (s *Service) GetLatestSpaceCache(ctx context.Context, source string) (*cache.SpaceCache, *rpc_errors.ServiceError) {
	spec := specification.NewGetLastNBySourceSpec(
		s.cacheRepo.GetIDFieldName(ctx),
		1,
		source,
	)

	result, err := s.cacheRepo.Find(ctx, spec)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, rpc_errors.NewNotFoundError(err)
		}
		return nil, rpc_errors.NewInternalError("Failed to find space cache data", err)
	}

	if len(result) == 0 {
		return nil, rpc_errors.NewNotFoundError(repository.ErrNotFound)
	}

	return &result[0], nil
}
