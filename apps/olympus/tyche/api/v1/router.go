package v1_tyche

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Echo) *echo.Echo {
	e.GET("/health", Health)
	return e
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}
