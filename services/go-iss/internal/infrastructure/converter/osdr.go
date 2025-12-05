package converter

import (
	"encoding/json"
	"go-iss/internal/infrastructure/repository/osdr"
	"go-iss/internal/rpc/models"
)

func ToOSDRItemResponse(item *osdr.OSDRItem) *models.OSDRItem {
	res := &models.OSDRItem{
		ID:         item.ID,
		DatasetID:  item.DatasetID,
		Title:      item.Title,
		Status:     item.Status,
		UpdatedAt:  item.UpdatedAt,
		InsertedAt: item.InsertedAt,
	}

	if err := json.Unmarshal(item.Raw, &res.Raw); err != nil {
		res.Raw = map[string]interface{}{}
	}

	return res
}

func ToOSDRListResponse(items []osdr.OSDRItem) *models.OSDRListResponse {
	result := make([]models.OSDRItem, 0, len(items))
	for i := range items {
		result = append(result, *ToOSDRItemResponse(&items[i]))
	}

	return &models.OSDRListResponse{
		Items: result,
	}
}
