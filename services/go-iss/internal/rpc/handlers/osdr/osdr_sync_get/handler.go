package osdr_sync_get

import (
	"context"
	"go-iss/internal/rpc"
	"go-iss/internal/rpc/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	FetchAndStoreOsdrService FetchAndStoreOsdrService
}

func New(service FetchAndStoreOsdrService) *Handler {
	return &Handler{
		FetchAndStoreOsdrService: service,
	}
}

func (h *Handler) OSDRSyncGet(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()

	traceID := rpc.GetTraceIDFromContext(requestCtx)
	requestCtx = context.WithValue(requestCtx, "trace_id", traceID)

	written, err := h.FetchAndStoreOsdrService.FetchAndStoreOSDR(requestCtx)
	if err != nil {
		errorResponse := models.NewErrorResponse(
			err.Code,
			err.Message,
			traceID,
		)
		return ctx.JSON(http.StatusOK, errorResponse)
	}

	successResponse := models.SuccessResponse(models.OSDRSyncResponse{Written: written})
	return ctx.JSON(http.StatusOK, successResponse)
}
