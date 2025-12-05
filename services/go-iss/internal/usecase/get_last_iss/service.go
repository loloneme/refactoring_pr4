package get_last_iss

import (
	"context"
	"errors"
	"go-iss/internal/domain/iss/specification"
	"go-iss/internal/infrastructure/repository"
	"go-iss/internal/infrastructure/repository/iss"
	rpc_errors "go-iss/internal/rpc/errors"
)

type Service struct {
	issRepo issRepo
}

func New(repo issRepo) *Service {
	return &Service{
		issRepo: repo,
	}
}

func (s *Service) GetLastISS(ctx context.Context) (*iss.ISSFetchLog, *rpc_errors.ServiceError) {
	spec := specification.NewGetLastNByIdDescSpec(
		s.issRepo.GetIDFieldName(ctx),
		1,
		[]string{"*"},
	)

	result, err := s.issRepo.Find(ctx, spec)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, rpc_errors.NewNotFoundError(err)
		}
		return nil, rpc_errors.NewInternalError("Failed to find ISS data:", err)
	}

	if len(result) == 0 {
		return nil, rpc_errors.NewNotFoundError(repository.ErrNotFound)
	}

	return &result[0], nil
}
