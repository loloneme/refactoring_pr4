package models

import "time"

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	TraceID string `json:"trace_id"`
}

type APIResponse struct {
	OK    bool           `json:"ok"`
	Error *ErrorResponse `json:"error,omitempty"`
	Data  interface{}    `json:"data,omitempty"`
}

type HealthResponse struct {
	Status string    `json:"status"`
	Now    time.Time `json:"now"`
}

func SuccessResponse(data interface{}) *APIResponse {
	return &APIResponse{
		OK:   true,
		Data: data,
	}
}

func NewErrorResponse(code, message, traceID string) *APIResponse {
	return &APIResponse{
		OK: false,
		Error: &ErrorResponse{
			Code:    code,
			Message: message,
			TraceID: traceID,
		},
	}
}
