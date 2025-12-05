package osdr_sync_get

import (
	"context"
	"go-iss/internal/rpc/errors"
)

type FetchAndStoreOsdrService interface {
	FetchAndStoreOSDR(ctx context.Context) (int, *errors.ServiceError)
}
