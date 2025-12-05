package json

import (
	"encoding/json"
	"fmt"
	"go-iss/internal/rpc/models"
	"time"

	"github.com/jmoiron/sqlx/types"
)

func ExtractString(v map[string]interface{}, keys []string) *string {
	for _, k := range keys {
		if x, ok := v[k]; ok {
			if s, ok := x.(string); ok && s != "" {
				return &s
			}
			if n, ok := x.(float64); ok {
				s := fmt.Sprintf("%.0f", n)
				return &s
			}
		}
	}
	return nil
}

func ExtractTime(v map[string]interface{}, keys []string) *time.Time {
	for _, k := range keys {
		if x, ok := v[k]; ok {
			if s, ok := x.(string); ok {
				formats := []string{
					time.RFC3339,
					"2006-01-02T15:04:05Z07:00",
					"2006-01-02 15:04:05",
					"2006-01-02T15:04:05",
				}
				for _, format := range formats {
					if t, err := time.Parse(format, s); err == nil {
						utc := t.UTC()
						return &utc
					}
				}
			}
			if n, ok := x.(float64); ok {
				t := time.Unix(int64(n), 0).UTC()
				return &t
			}
		}
	}
	return nil
}

func ExtractItems(jsonData interface{}) []interface{} {
	switch v := jsonData.(type) {
	case []interface{}:
		return v
	case map[string]interface{}:
		if items, ok := v["items"].([]interface{}); ok {
			return items
		}
		if results, ok := v["results"].([]interface{}); ok {
			return results
		}
		return []interface{}{v}
	default:
		return []interface{}{v}
	}
}

func MarshalPayloadToCache(item *models.SpaceCacheItem, payload types.JSONText) models.SpaceCacheItem {
	if err := json.Unmarshal(payload, &item.Payload); err != nil {
		item.Payload = map[string]interface{}{}
	}

	return *item
}
