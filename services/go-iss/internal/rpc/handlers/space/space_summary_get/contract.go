package space_summary_get

import (
	"context"
	"go-iss/internal/rpc/errors"
	"go-iss/internal/rpc/models"
)

type GetSpaceSummaryService interface {
	GetSpaceSummary(ctx context.Context) (*models.SpaceSummary, *errors.ServiceError)
}
