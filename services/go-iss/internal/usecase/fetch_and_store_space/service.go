package fetch_and_store_space

import (
	"context"
	"encoding/json"
	"go-iss/internal/infrastructure/repository/cache"
	rpc_errors "go-iss/internal/rpc/errors"
	time_util "go-iss/internal/usecase/utils/time"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx/types"
)

type Service struct {
	cacheRepo    cacheRepo
	nasaClient   nasaClient
	spacexClient spacexClient
}

func New(repo cacheRepo, nasa nasaClient, spacex spacexClient) *Service {
	return &Service{
		cacheRepo:    repo,
		nasaClient:   nasa,
		spacexClient: spacex,
	}
}

func (s *Service) RefreshSpace(ctx context.Context, sources []string) ([]string, *rpc_errors.ServiceError) {
	var refreshed []string

	for _, source := range sources {
		source = strings.TrimSpace(strings.ToLower(source))

		var jsonData interface{}
		var err error

		switch source {
		case "apod":
			jsonData, err = s.nasaClient.FetchAPOD(ctx)
		case "neo":
			from, to := time_util.LastDays(2)
			jsonData, err = s.nasaClient.FetchNEOFeed(ctx, from, to)
		case "flr":
			from, to := time_util.LastDays(5)
			jsonData, err = s.nasaClient.FetchDONKIFLR(ctx, from, to)
		case "cme":
			from, to := time_util.LastDays(5)
			jsonData, err = s.nasaClient.FetchDONKICME(ctx, from, to)
		case "spacex":
			jsonData, err = s.spacexClient.FetchNextLaunch(ctx)
		default:
			continue
		}

		if err != nil {
			log.Printf("Error fetching %s: %v\n", source, err.Error())
			continue
		}

		payloadJSON, err := json.Marshal(jsonData)
		if err != nil {
			log.Printf("Error marshalling JSON from %s: %v\n", source, err.Error())
			continue
		}

		entity := cache.SpaceCache{
			Source:    source,
			FetchedAt: time.Now().UTC(),
			Payload:   types.JSONText(payloadJSON),
		}

		if err := s.cacheRepo.Save(ctx, entity); err != nil {
			log.Printf("Error saving from %s: %v\n", source, err.Error())
			continue
		}

		refreshed = append(refreshed, source)
	}

	return refreshed, nil
}
