package iss_trend_get

import (
	"context"
	"go-iss/internal/rpc"
	"go-iss/internal/rpc/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	GetIssTrendService GetIssTrendService
}

func New(service GetIssTrendService) *Handler {
	return &Handler{
		GetIssTrendService: service,
	}
}

func (h *Handler) ISSTrendGet(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()

	traceID := rpc.GetTraceIDFromContext(requestCtx)
	requestCtx = context.WithValue(requestCtx, "trace_id", traceID)

	result, err := h.GetIssTrendService.GetISSTrend(requestCtx)
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
