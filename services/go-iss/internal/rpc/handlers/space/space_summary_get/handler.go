package space_summary_get

import (
	"context"
	"go-iss/internal/rpc"
	"go-iss/internal/rpc/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	GetSpaceSummaryService GetSpaceSummaryService
}

func New(service GetSpaceSummaryService) *Handler {
	return &Handler{
		GetSpaceSummaryService: service,
	}
}

func (h *Handler) SpaceSummaryGet(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()

	traceID := rpc.GetTraceIDFromContext(requestCtx)
	requestCtx = context.WithValue(requestCtx, "trace_id", traceID)

	result, err := h.GetSpaceSummaryService.GetSpaceSummary(requestCtx)
	if err != nil {
		errorResponse := models.NewErrorResponse(
			err.Code,
			err.Message,
			traceID,
		)
		return ctx.JSON(http.StatusOK, errorResponse)
	}

	successResponse := models.SuccessResponse(result)
	return ctx.JSON(http.StatusOK, successResponse)
}
