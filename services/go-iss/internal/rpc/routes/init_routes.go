package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRoutes(e *echo.Echo, handlers Handlers) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", handlers.HandleHealth)
	e.GET("/iss/last", handlers.HandleLastISS)
	e.GET("/iss/fetch", handlers.HandleTriggerISS)
	e.GET("/iss/trend", handlers.HandleISSTrend)
	e.GET("/osdr/sync", handlers.HandleOSDRSync)
	e.GET("/osdr/list", handlers.HandleOSDRList)
	e.GET("/space/:src/latest", handlers.HandleSpaceLatest)
	e.GET("/space/refresh", handlers.HandleSpaceRefresh)
	e.GET("/space/summary", handlers.HandleSpaceSummary)
}

type Handlers struct {
	HandleHealth       echo.HandlerFunc
	HandleLastISS      echo.HandlerFunc
	HandleTriggerISS   echo.HandlerFunc
	HandleISSTrend     echo.HandlerFunc
	HandleOSDRSync     echo.HandlerFunc
	HandleOSDRList     echo.HandlerFunc
	HandleSpaceLatest  echo.HandlerFunc
	HandleSpaceRefresh echo.HandlerFunc
	HandleSpaceSummary echo.HandlerFunc
}
