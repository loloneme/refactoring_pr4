package iss_last_get

import (
	"context"
	"go-iss/internal/infrastructure/converter"
	"go-iss/internal/rpc"
	"net/http"

	"go-iss/internal/rpc/models"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	GetLastIssService GetLastIssService
}

func New(service GetLastIssService) *Handler {
	return &Handler{
		GetLastIssService: service,
	}
}

func (h *Handler) ISSLastGet(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()

	traceID := rpc.GetTraceIDFromContext(requestCtx)
	requestCtx = context.WithValue(requestCtx, "trace_id", traceID)

	result, err := h.GetLastIssService.GetLastISS(requestCtx)
	if err != nil {
		errorResponse := models.NewErrorResponse(
			err.Code,
			err.Message,
			traceID,
		)
		return ctx.JSON(http.StatusOK, errorResponse)
	}

	successResponse := models.SuccessResponse(converter.ToISSJsonResponse(result))
	return ctx.JSON(http.StatusOK, successResponse)
}
