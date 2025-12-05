package space_refresh_get

import (
	"context"
	"go-iss/internal/rpc"
	"go-iss/internal/rpc/models"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	RefreshSpaceService RefreshSpaceService
}

func New(service RefreshSpaceService) *Handler {
	return &Handler{
		RefreshSpaceService: service,
	}
}

func (h *Handler) SpaceRefreshGet(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()

	traceID := rpc.GetTraceIDFromContext(requestCtx)
	requestCtx = context.WithValue(requestCtx, "trace_id", traceID)

	srcParam := ctx.QueryParam("src")
	if srcParam == "" {
		srcParam = "apod,neo,flr,cme,spacex"
	}

	sources := strings.Split(srcParam, ",")
	for i := range sources {
		sources[i] = strings.TrimSpace(sources[i])
	}

	refreshed, err := h.RefreshSpaceService.RefreshSpace(requestCtx, sources)
	if err != nil {
		errorResponse := models.NewErrorResponse(
			err.Code,
			err.Message,
			traceID,
		)
		return ctx.JSON(http.StatusOK, errorResponse)
	}

	response := models.SpaceRefreshResponse{
		Refreshed: refreshed,
	}
	successResponse := models.SuccessResponse(response)
	return ctx.JSON(http.StatusOK, successResponse)
}
