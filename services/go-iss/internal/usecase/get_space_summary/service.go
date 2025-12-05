package get_space_summary

import (
	"context"
	iss_specification "go-iss/internal/domain/iss/specification"
	space_cache_specification "go-iss/internal/domain/space_cache/specification"
	rpc_errors "go-iss/internal/rpc/errors"
	"go-iss/internal/rpc/models"
	json_util "go-iss/internal/usecase/utils/json"
	"log"
)

type Service struct {
	cacheRepo cacheRepo
	issRepo   issRepo
	osdrRepo  osdrRepo
}

func New(cache cacheRepo, iss issRepo, osdr osdrRepo) *Service {
	return &Service{
		cacheRepo: cache,
		issRepo:   iss,
		osdrRepo:  osdr,
	}
}

func (s *Service) GetSpaceSummary(ctx context.Context) (*models.SpaceSummary, *rpc_errors.ServiceError) {
	summary := &models.SpaceSummary{}

	sources := []string{"apod", "neo", "flr", "cme", "spacex"}
	for _, source := range sources {
		item := s.getLatestCacheItem(ctx, source)
		switch source {
		case "apod":
			summary.APOD = item
		case "neo":
			summary.NEO = item
		case "flr":
			summary.FLR = item
		case "cme":
			summary.CME = item
		case "spacex":
			summary.SpaceX = item
		}
	}

	issItem := s.getLastISSItem(ctx)
	summary.ISS = issItem

	osdrCount, err := s.osdrRepo.Count(ctx)
	if err != nil {
		log.Printf("Error counting OSDR items: %v", err)
	}
	summary.OSDRCount = osdrCount

	return summary, nil
}

func (s *Service) getLatestCacheItem(ctx context.Context, source string) models.SpaceCacheItem {
	spec := space_cache_specification.NewGetLastNBySourceSpec(
		s.cacheRepo.GetIDFieldName(ctx),
		1,
		source,
	)

	result, err := s.cacheRepo.Find(ctx, spec)
	if err != nil || len(result) == 0 {
		return models.SpaceCacheItem{}
	}

	cacheItem := result[0]

	return json_util.MarshalPayloadToCache(&models.SpaceCacheItem{
		At: &cacheItem.FetchedAt,
	}, cacheItem.Payload)
}

func (s *Service) getLastISSItem(ctx context.Context) models.SpaceCacheItem {
	spec := iss_specification.NewGetLastNByIdDescSpec(
		s.issRepo.GetIDFieldName(ctx),
		1,
		[]string{"fetched_at", "payload"},
	)

	result, err := s.issRepo.Find(ctx, spec)
	if err != nil || len(result) == 0 {
		return models.SpaceCacheItem{}
	}

	issItem := result[0]

	return json_util.MarshalPayloadToCache(&models.SpaceCacheItem{
		At: &issItem.FetchedAt,
	}, issItem.Payload)
}
