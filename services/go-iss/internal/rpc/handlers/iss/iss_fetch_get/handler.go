package iss_fetch_get

import (
	"context"
	"go-iss/internal/infrastructure/converter"
	"go-iss/internal/rpc"
	"net/http"

	"go-iss/internal/rpc/models"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	FetchAndStoreIssService FetchAndStoreIssService
}

func New(FetchAndStoreIssService FetchAndStoreIssService) *Handler {
	return &Handler{
		FetchAndStoreIssService: FetchAndStoreIssService,
	}
}

func (h *Handler) ISSFetchGet(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()

	traceID := rpc.GetTraceIDFromContext(requestCtx)
	requestCtx = context.WithValue(requestCtx, "trace_id", traceID)

	result, err := h.FetchAndStoreIssService.FetchAndStoreISS(requestCtx)
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
