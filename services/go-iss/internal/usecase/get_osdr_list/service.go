package get_osdr_list

import (
	"context"
	"errors"
	"go-iss/internal/domain/osdr/specification"
	"go-iss/internal/infrastructure/repository"
	"go-iss/internal/infrastructure/repository/osdr"
	rpc_errors "go-iss/internal/rpc/errors"
)

type Service struct {
	osdrRepo osdrRepo
}

func New(repo osdrRepo) *Service {
	return &Service{
		osdrRepo: repo,
	}
}

func (s *Service) GetOSDRList(ctx context.Context, limit int) ([]osdr.OSDRItem, *rpc_errors.ServiceError) {
	if limit <= 0 {
		limit = 20
	}

	spec := specification.NewGetLastNByIdDescSpec(limit)

	items, err := s.osdrRepo.Find(ctx, spec)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return []osdr.OSDRItem{}, nil
		}
		return nil, rpc_errors.NewInternalError("Failed to find OSDR items", err)
	}

	return items, nil
}
