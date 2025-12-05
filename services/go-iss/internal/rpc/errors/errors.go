package errors

import (
	"fmt"
	"net/http"
)

type ServiceError struct {
	Code    string
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}

const (
	ErrorCodeUpstream403     string = "UPSTREAM_403"
	ErrorCodeUpstream404     string = "UPSTREAM_404"
	ErrorCodeUpstream500     string = "UPSTREAM_500"
	ErrorCodeUpstreamTimeout string = "UPSTREAM_TIMEOUT"
	ErrorCodeUpstreamGeneric string = "UPSTREAM_ERROR"
	ErrorCodeNotFoundError   string = "RESOURCE_NOT_FOUND"
	ErrorCodeInternalError   string = "INTERNAL_ERROR"
)

// NewUpstreamError создает ошибку upstream на основе HTTP статуса
func NewUpstreamError(statusCode int, err error) *ServiceError {
	var code string
	var message string

	switch statusCode {
	case http.StatusForbidden:
		code = ErrorCodeUpstream403
		message = "Upstream service returned 403 Forbidden"
	case http.StatusNotFound:
		code = ErrorCodeUpstream404
		message = "Upstream service returned 404 Not Found"
	case http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		code = ErrorCodeUpstream500
		message = "Upstream service error"
	default:
		code = ErrorCodeUpstreamGeneric
		message = "Upstream service error"
	}

	return &ServiceError{
		Code:    code,
		Message: fmt.Sprintf("%s: %v", message, err),
	}
}

func NewUpstreamTimeoutError(err error) *ServiceError {
	return &ServiceError{
		Code:    ErrorCodeUpstreamTimeout,
		Message: fmt.Sprintf("Request timeout: %v", err),
	}
}

func NewUpstreamGenericError(err error) *ServiceError {
	return &ServiceError{
		Code:    ErrorCodeUpstreamGeneric,
		Message: fmt.Sprintf("Failed to fetch data from upstream service: %v", err),
	}
}

func NewInternalError(message string, err error) *ServiceError {
	return &ServiceError{
		Code:    ErrorCodeInternalError,
		Message: fmt.Sprintf("%s: %v", message, err),
	}
}

func NewNotFoundError(err error) *ServiceError {
	return &ServiceError{
		Code:    ErrorCodeNotFoundError,
		Message: fmt.Sprintf("%v", err),
	}
}
