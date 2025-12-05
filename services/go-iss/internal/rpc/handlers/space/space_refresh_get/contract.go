package space_refresh_get

import (
	"context"
	"go-iss/internal/rpc/errors"
)

type RefreshSpaceService interface {
	RefreshSpace(ctx context.Context, sources []string) ([]string, *errors.ServiceError)
}
