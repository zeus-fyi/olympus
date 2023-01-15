package v1hestia

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Routes(e *echo.Echo) *echo.Echo {
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Routes
	e.GET("/health", Health)

	// TODO add users/orgs
	// TODO add validator service groups

	return e
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}
