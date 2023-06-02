package v1_iris

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Echo) *echo.Echo {
	// Routes
	e.POST("/", Proxy)
	e.GET("/health", Health)
	return e
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}
