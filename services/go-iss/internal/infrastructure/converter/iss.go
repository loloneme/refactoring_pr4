package converter

import (
	"encoding/json"
	"go-iss/internal/infrastructure/repository/iss"
	"go-iss/internal/rpc/models"
)

func ToISSJsonResponse(log *iss.ISSFetchLog) *models.ISSFetchLog {
	res := &models.ISSFetchLog{
		ID:        log.ID,
		FetchedAt: log.FetchedAt,
		SourceURL: log.SourceURL,
	}

	if err := json.Unmarshal(log.Payload, &res.Payload); err != nil {
		res.Payload = map[string]interface{}{}
	}

	return res
}
