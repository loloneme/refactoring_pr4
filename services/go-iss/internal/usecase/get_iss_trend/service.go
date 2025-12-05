package get_iss_trend

import (
	"context"
	"encoding/json"
	"errors"
	"go-iss/internal/domain/iss/specification"
	"go-iss/internal/infrastructure/repository"
	"go-iss/internal/infrastructure/repository/iss"
	rpc_errors "go-iss/internal/rpc/errors"
	"go-iss/internal/rpc/models"
	"go-iss/internal/usecase/utils/distance"
	"go-iss/internal/usecase/utils/float"
	"time"
)

type Service struct {
	issRepo issRepo
}

func New(repo issRepo) *Service {
	return &Service{
		issRepo: repo,
	}
}

func (s *Service) GetISSTrend(ctx context.Context) (*models.Trend, *rpc_errors.ServiceError) {
	records, err := s.getLastTwoRecords(ctx)
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		return &models.Trend{}, nil
	}

	payload1, payload2, err := s.parsePayloads(records[1], records[0])
	if err != nil {
		return nil, err
	}

	coords := s.extractCoordinates(payload1, payload2)
	movement := s.calculateMovement(coords)

	trend := s.buildTrend(records[1].FetchedAt, records[0].FetchedAt, coords, movement)

	return trend, nil
}

func (s *Service) getLastTwoRecords(ctx context.Context) ([]iss.ISSFetchLog, *rpc_errors.ServiceError) {
	idField := s.issRepo.GetIDFieldName(ctx)

	spec := specification.NewGetLastNByIdDescSpec(
		idField,
		2,
		[]string{"fetched_at", "payload"},
	)

	records, err := s.issRepo.Find(ctx, spec)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return []iss.ISSFetchLog{}, nil
		}
		return nil, rpc_errors.NewInternalError("Failed to find ISS data", err)
	}

	return records, nil
}

func (s *Service) parsePayloads(record1, record2 iss.ISSFetchLog) (map[string]interface{}, map[string]interface{}, *rpc_errors.ServiceError) {
	var payload1, payload2 map[string]interface{}

	if err := json.Unmarshal(record1.Payload, &payload1); err != nil {
		return nil, nil, rpc_errors.NewInternalError("Failed to unmarshal payload", err)
	}

	if err := json.Unmarshal(record2.Payload, &payload2); err != nil {
		return nil, nil, rpc_errors.NewInternalError("Failed to unmarshal payload", err)
	}

	return payload1, payload2, nil
}

type coordinates struct {
	lat1, lon1, lat2, lon2, velocity *float64
}

func (s *Service) extractCoordinates(payload1, payload2 map[string]interface{}) coordinates {
	return coordinates{
		lat1:     float.ExtractFloat(payload1, "latitude"),
		lon1:     float.ExtractFloat(payload1, "longitude"),
		lat2:     float.ExtractFloat(payload2, "latitude"),
		lon2:     float.ExtractFloat(payload2, "longitude"),
		velocity: float.ExtractFloat(payload2, "velocity"),
	}
}

type movementResult struct {
	deltaKm  float64
	movement bool
}

func (s *Service) calculateMovement(coords coordinates) movementResult {
	if coords.lat1 == nil || coords.lon1 == nil || coords.lat2 == nil || coords.lon2 == nil {
		return movementResult{deltaKm: 0.0, movement: false}
	}

	deltaKm := distance.HaversineKm(*coords.lat1, *coords.lon1, *coords.lat2, *coords.lon2)
	movement := deltaKm > 0.1

	return movementResult{deltaKm: deltaKm, movement: movement}
}

func (s *Service) buildTrend(fromTime, toTime time.Time, coords coordinates, movement movementResult) *models.Trend {
	dtSec := toTime.Sub(fromTime).Seconds()

	trend := models.Trend{
		Movement: movement.movement,
		DeltaKm:  movement.deltaKm,
		DtSec:    dtSec,
		FromTime: &fromTime,
		ToTime:   &toTime,
	}

	if coords.velocity != nil {
		trend.VelocityKmh = coords.velocity
	}
	if coords.lat1 != nil {
		trend.FromLat = coords.lat1
	}
	if coords.lon1 != nil {
		trend.FromLon = coords.lon1
	}
	if coords.lat2 != nil {
		trend.ToLat = coords.lat2
	}
	if coords.lon2 != nil {
		trend.ToLon = coords.lon2
	}

	return &trend
}
