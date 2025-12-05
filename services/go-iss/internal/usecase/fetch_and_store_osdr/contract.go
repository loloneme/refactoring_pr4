package fetch_and_store_osdr

import (
	"context"
	"go-iss/internal/infrastructure/repository/osdr"
)

type osdrRepo interface {
	SaveOSDR(ctx context.Context, entity osdr.OSDRItem) error
}

type nasaClient interface {
	FetchOSDR(ctx context.Context, osdrURL string) (interface{}, error)
	GetOSDRURL() string
}
