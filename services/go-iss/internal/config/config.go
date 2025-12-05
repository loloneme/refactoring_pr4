package config

import (
	"context"
	"os"
	"strconv"

	iss_client "go-iss/internal/infrastructure/clients/iss"
	"go-iss/internal/infrastructure/clients/nasa"
	"go-iss/internal/infrastructure/clients/spacex"
	cache_repo "go-iss/internal/infrastructure/repository/cache"
	iss_repo "go-iss/internal/infrastructure/repository/iss"
	osdr_repo "go-iss/internal/infrastructure/repository/osdr"

	"github.com/jmoiron/sqlx"
)

type RetryConfig struct {
	EveryOSDR   uint64
	EveryISS    uint64
	EveryAPOD   uint64
	EveryNEO    uint64
	EveryDONKI  uint64
	EverySpaceX uint64
}

type BackgroundConfig struct {
	RetryTimes   RetryConfig
	ISSClient    *iss_client.Client
	NasaClient   *nasa.Client
	SpaceXClient *spacex.Client
	ISSRepo      *iss_repo.Repository
	OSDRRepo     *osdr_repo.Repository
	CacheRepo    *cache_repo.Repository
}

func SetupBackgroundConfig(ctx context.Context, db *sqlx.DB, issClient *iss_client.Client, nasaClient *nasa.Client, spacexClient *spacex.Client) (*BackgroundConfig, error) {
	retryConfig := RetryConfig{
		EveryOSDR:   envU64("FETCH_EVERY_SECONDS", 600),
		EveryISS:    envU64("ISS_EVERY_SECONDS", 120),
		EveryAPOD:   envU64("APOD_EVERY_SECONDS", 43200), // 12ч
		EveryNEO:    envU64("NEO_EVERY_SECONDS", 7200),   // 2ч
		EveryDONKI:  envU64("DONKI_EVERY_SECONDS", 3600), // 1ч
		EverySpaceX: envU64("SPACEX_EVERY_SECONDS", 3600),
	}

	issRepo := iss_repo.NewRepository(db)
	osdrRepo := osdr_repo.NewRepository(db)
	cacheRepo := cache_repo.NewRepository(db)

	cfg := &BackgroundConfig{
		RetryTimes:   retryConfig,
		ISSClient:    issClient,
		NasaClient:   nasaClient,
		SpaceXClient: spacexClient,
		ISSRepo:      issRepo,
		OSDRRepo:     osdrRepo,
		CacheRepo:    cacheRepo,
	}

	return cfg, nil
}

func envU64(key string, defaultValue uint64) uint64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseUint(value, 10, 64); err == nil {
			return parsed
		}
	}
	return defaultValue
}
