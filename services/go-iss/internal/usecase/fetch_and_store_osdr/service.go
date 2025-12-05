package fetch_and_store_osdr

import (
	"context"
	"encoding/json"
	"fmt"
	"go-iss/internal/infrastructure/repository/osdr"
	rpc_errors "go-iss/internal/rpc/errors"
	"strings"

	"github.com/jmoiron/sqlx/types"
)

type Service struct {
	osdrRepo   osdrRepo
	nasaClient nasaClient
}

func New(repo osdrRepo, client nasaClient) *Service {
	return &Service{
		osdrRepo:   repo,
		nasaClient: client,
	}
}

func (s *Service) FetchAndStoreOSDR(ctx context.Context) (int, *rpc_errors.ServiceError) {
	osdrURL := s.nasaClient.GetOSDRURL()

	jsonData, err := s.nasaClient.FetchOSDR(ctx, osdrURL)
	if err != nil {
		errMsg := strings.ToLower(err.Error())
		if strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "deadline exceeded") {
			return 0, rpc_errors.NewUpstreamTimeoutError(err)
		}
		return 0, rpc_errors.NewUpstreamGenericError(err)
	}

	dataMap, ok := jsonData.(map[string]interface{})
	if !ok {
		return 0, rpc_errors.NewUpstreamGenericError(fmt.Errorf("unexpected response format"))
	}

	//items := json_utils.ExtractItems(jsonData)

	written := 0
	for datasetKey, value := range dataMap {
		valueMap, ok := value.(map[string]interface{})
		if !ok {
			continue
		}

		datasetID := datasetKey
		title := datasetKey
		//title := json_utils.ExtractString(valueMap, []string{"title", "name", "label"})
		//status := json_utils.ExtractString(valueMap, []string{"status", "state", "lifecycle"})
		//updatedAt := json_utils.ExtractTime(valueMap, []string{"updated", "updated_at", "modified", "lastUpdated", "timestamp"})

		itemJSON, _ := json.Marshal(valueMap)

		entity := osdr.OSDRItem{
			DatasetID: &datasetID,
			Title:     &title,
			//Status:     status,
			//UpdatedAt:  updatedAt,
			//InsertedAt: now,
			Raw: types.JSONText(itemJSON),
		}

		if err := s.osdrRepo.SaveOSDR(ctx, entity); err != nil {
			continue
		}

		written++
	}

	return written, nil
}
