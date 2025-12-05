package handlers

import (
	"go-iss/internal/rpc/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func HealthCheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, models.HealthResponse{
		Status: "ok",
		Now:    time.Now().UTC(),
	})
}
