package iss_last_get

import (
	"context"
	"go-iss/internal/infrastructure/repository/iss"
	rpc_errors "go-iss/internal/rpc/errors"
)

type GetLastIssService interface {
	GetLastISS(ctx context.Context) (*iss.ISSFetchLog, *rpc_errors.ServiceError)
}
