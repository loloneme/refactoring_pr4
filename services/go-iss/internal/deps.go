package internal

import (
	"context"
	"go-iss/internal/infrastructure/clients/iss"
	"go-iss/internal/infrastructure/clients/nasa"
	"go-iss/internal/infrastructure/clients/spacex"
	"go-iss/internal/infrastructure/postgres"

	"github.com/jmoiron/sqlx"
)

func NewDatabaseConnection(ctx context.Context) (*sqlx.DB, error) {
	return postgres.NewFromConfig(ctx)
}

func NewNASAClient() (*nasa.Client, error) {
	return nasa.NewFromConfig()
}

func NewISSClient() (*iss.Client, error) {
	return iss.NewFromConfig()
}

func NewSpaceXClient() (*spacex.Client, error) {
	return spacex.NewFromConfig()
}
