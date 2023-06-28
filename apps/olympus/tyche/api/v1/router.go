package v1_tyche

import (
	"net/http"

	"github.com/labstack/echo/v4"
	tyche_metrics "github.com/zeus-fyi/olympus/tyche/metrics"
)

func Routes(e *echo.Echo) *echo.Echo {
	e.GET("/health", Health)
	e.POST(txProcessorRoute, TxProcessingRequestHandler)
	e.GET("/metrics", tyche_metrics.MetricsRequestHandler)
	return e
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}
