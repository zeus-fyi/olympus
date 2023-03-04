package hestia_web_router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func WebRoutes(e *echo.Echo) *echo.Echo {
	// Routes

	e.GET("/login", Login)
	return e
}

func Login(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}
