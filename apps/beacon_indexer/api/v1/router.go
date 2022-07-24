package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Routes(e *echo.Echo) *echo.Echo {
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Routes
	e.GET("/health", Health)
	e.POST("/admin", HandleAdminConfigRequest)
	e.GET("/admin", HandleAdminGetRequest)

	e.GET("/debug", HandleDebugRequest)
	return e
}
