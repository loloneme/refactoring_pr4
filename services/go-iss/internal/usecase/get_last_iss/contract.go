package get_last_iss

import (
	"context"
	"go-iss/internal/infrastructure/repository"
	"go-iss/internal/infrastructure/repository/iss"
)

type issRepo interface {
	Find(ctx context.Context, spec repository.FindSpecification) ([]iss.ISSFetchLog, error)
	GetIDFieldName(ctx context.Context) string
}
