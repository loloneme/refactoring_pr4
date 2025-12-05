package iss_fetch_get

import (
	"context"
	"go-iss/internal/infrastructure/repository/iss"
	"go-iss/internal/rpc/errors"
)

type FetchAndStoreIssService interface {
	FetchAndStoreISS(ctx context.Context) (*iss.ISSFetchLog, *errors.ServiceError)
}
