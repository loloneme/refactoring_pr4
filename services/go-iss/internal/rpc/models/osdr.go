package models

import "time"

type OSDRItem struct {
	ID         int64                  `json:"id"`
	DatasetID  *string                `json:"dataset_id,omitempty"`
	Title      *string                `json:"title,omitempty"`
	Status     *string                `json:"status,omitempty"`
	UpdatedAt  *time.Time             `json:"updated_at,omitempty"`
	InsertedAt time.Time              `json:"inserted_at"`
	Raw        map[string]interface{} `json:"raw"`
}

type OSDRListResponse struct {
	Items []OSDRItem `json:"items"`
}

type OSDRSyncResponse struct {
	Written int `json:"written"`
}
