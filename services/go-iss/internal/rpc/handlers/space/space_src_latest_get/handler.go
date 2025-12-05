package space_src_latest_get

import (
	"context"
	"go-iss/internal/infrastructure/converter"
	"go-iss/internal/rpc"
	"go-iss/internal/rpc/errors"
	"go-iss/internal/rpc/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	GetLatestSpaceCacheService GetLatestSpaceCacheService
}

func New(service GetLatestSpaceCacheService) *Handler {
	return &Handler{
		GetLatestSpaceCacheService: service,
	}
}

func (h *Handler) SpaceSrcLatestGet(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()

	traceID := rpc.GetTraceIDFromContext(requestCtx)
	requestCtx = context.WithValue(requestCtx, "trace_id", traceID)

	source := ctx.Param("src")

	result, err := h.GetLatestSpaceCacheService.GetLatestSpaceCache(requestCtx, source)
	if err != nil {
		if err.Code == errors.ErrorCodeNotFoundError {
			noDataResponse := models.SpaceLatestNoDataResponse{
				Source:  source,
				Message: "no data",
			}
			successResponse := models.SuccessResponse(noDataResponse)
			return ctx.JSON(http.StatusOK, successResponse)
		}
		errorResponse := models.NewErrorResponse(
			err.Code,
			err.Message,
			traceID,
		)
		return ctx.JSON(http.StatusOK, errorResponse)
	}

	successResponse := models.SuccessResponse(converter.ToSpaceLatestResponse(result))
	return ctx.JSON(http.StatusOK, successResponse)
}
