package iss_trend_get

import (
	"context"
	"go-iss/internal/rpc/errors"
	"go-iss/internal/rpc/models"
)

type GetIssTrendService interface {
	GetISSTrend(ctx context.Context) (*models.Trend, *errors.ServiceError)
}
