package fetch_and_store_iss

import (
	"context"
	"encoding/json"
	"errors"
	"go-iss/internal/domain/iss/specification"
	"go-iss/internal/infrastructure/repository"
	"go-iss/internal/infrastructure/repository/iss"
	rpc_errors "go-iss/internal/rpc/errors"
	"strings"

	"github.com/jmoiron/sqlx/types"
)

type Service struct {
	issRepo   issRepo
	issClient issClient
}

func New(repo issRepo, client issClient) *Service {
	return &Service{
		issRepo:   repo,
		issClient: client,
	}
}

func (s *Service) FetchAndStoreISS(ctx context.Context) (*iss.ISSFetchLog, *rpc_errors.ServiceError) {
	jsonData, err, code := s.issClient.FetchISS(ctx)
	if err != nil {
		if code == 0 {
			errMsg := strings.ToLower(err.Error())
			if strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "deadline exceeded") {
				return nil, rpc_errors.NewUpstreamTimeoutError(err)
			}
			return nil, rpc_errors.NewUpstreamGenericError(err)
		}
		return nil, rpc_errors.NewUpstreamError(code, err)
	}

	payloadJSON, err := json.Marshal(jsonData)
	if err != nil {
		return nil, rpc_errors.NewInternalError("Failed to marshal ISS data", err)
	}

	entity := iss.ISSFetchLog{
		SourceURL: s.issClient.GetSourceURL(),
		Payload:   types.JSONText(payloadJSON),
	}

	err = s.issRepo.Save(ctx, entity)
	if err != nil {
		return nil, rpc_errors.NewInternalError("Failed to save ISS data", err)
	}

	idField := s.issRepo.GetIDFieldName(ctx)

	spec := specification.NewGetLastNByIdDescSpec(
		idField,
		1,
		[]string{
			idField, "fetched_at", "source_url", "payload",
		},
	)

	createdLog, err := s.issRepo.Find(ctx, spec)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, rpc_errors.NewNotFoundError(err)
		}
		return nil, rpc_errors.NewInternalError("Failed to find ISS data", err)
	}

	return &createdLog[0], nil
}
