package osdr_list_get

import (
	"context"
	"go-iss/internal/infrastructure/repository/osdr"
	"go-iss/internal/rpc/errors"
)

type GetOsdrListService interface {
	GetOSDRList(ctx context.Context, limit int) ([]osdr.OSDRItem, *errors.ServiceError)
}
