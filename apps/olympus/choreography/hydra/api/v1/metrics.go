package v1_hydra_choreography

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	hyrda_choreography_metrics "github.com/zeus-fyi/olympus/choreography/hydra/metrics"
)

func MetricsHandler(c echo.Context) error {
	return MetricsRequest(c)
}

func MetricsRequest(c echo.Context) error {
	promHandler := promhttp.HandlerFor(hyrda_choreography_metrics.HydraChoreographyMetrics, promhttp.HandlerOpts{})
	return echo.WrapHandler(promHandler)(c)
}
