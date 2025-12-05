package get_osdr_list

import (
	"context"
	"go-iss/internal/infrastructure/repository"
	"go-iss/internal/infrastructure/repository/osdr"
)

type osdrRepo interface {
	Find(ctx context.Context, spec repository.FindSpecification) ([]osdr.OSDRItem, error)
}
