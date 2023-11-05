package v1_promql

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Echo) *echo.Echo {
	// Routes
	e.GET("/health", Health)
	e.GET("/v1/promql/top/tokens", PromQLRequestHandler("top_k_tokens"))
	return e
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}
