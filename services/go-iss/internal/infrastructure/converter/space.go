package converter

import (
	"encoding/json"
	"go-iss/internal/infrastructure/repository/cache"
	"go-iss/internal/rpc/models"
)

func ToSpaceLatestResponse(model *cache.SpaceCache) *models.SpaceLatestResponse {
	res := &models.SpaceLatestResponse{
		Source:    model.Source,
		FetchedAt: model.FetchedAt,
	}
	if err := json.Unmarshal(model.Payload, &res.Payload); err != nil {
		res.Payload = map[string]interface{}{}
	}

	return res
}
