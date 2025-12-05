package osdr_list_get

import (
	"context"
	"go-iss/internal/infrastructure/converter"
	"go-iss/internal/rpc"
	"go-iss/internal/rpc/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	GetOsdrListService GetOsdrListService
}

func New(service GetOsdrListService) *Handler {
	return &Handler{
		GetOsdrListService: service,
	}
}

func (h *Handler) OSDRListGet(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()

	traceID := rpc.GetTraceIDFromContext(requestCtx)
	requestCtx = context.WithValue(requestCtx, "trace_id", traceID)

	limitStr := ctx.QueryParam("limit")
	limit := 20
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	result, err := h.GetOsdrListService.GetOSDRList(requestCtx, limit)
	if err != nil {
		errorResponse := models.NewErrorResponse(
			err.Code,
			err.Message,
			traceID,
		)
		return ctx.JSON(http.StatusOK, errorResponse)
	}

	successResponse := models.SuccessResponse(converter.ToOSDRListResponse(result))
	return ctx.JSON(http.StatusOK, successResponse)
}
