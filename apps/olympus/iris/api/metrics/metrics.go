package iris_metrics

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	apollo_metrics "github.com/zeus-fyi/olympus/pkg/apollo/metrics"
)

var IrisMetrics *prometheus.Registry

func MetricsRequestHandler(c echo.Context) error {
	promHandler := promhttp.HandlerFor(IrisMetrics, promhttp.HandlerOpts{Registry: IrisMetrics})
	return echo.WrapHandler(promHandler)(c)
}

func InitIrisMetrics() {
	IrisMetrics = apollo_metrics.InitMetrics()
}
