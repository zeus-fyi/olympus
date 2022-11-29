package v1_hera

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Echo) *echo.Echo {
	// Routes
	e.GET("/health", Health)

	e.POST("/v1beta/ethereum/validator/deposits", Health)

	return e
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}
