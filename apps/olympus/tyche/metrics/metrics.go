package tyche_metrics

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics_trading "github.com/zeus-fyi/olympus/pkg/apollo/ethereum/mev/trading"
	apollo_metrics "github.com/zeus-fyi/olympus/pkg/apollo/metrics"
)

var (
	TycheMetrics *prometheus.Registry
	TradeMetrics metrics_trading.TradingMetrics
)

func MetricsRequestHandler(c echo.Context) error {
	promHandler := promhttp.HandlerFor(TycheMetrics, promhttp.HandlerOpts{Registry: TycheMetrics})
	return echo.WrapHandler(promHandler)(c)
}

func InitTycheMetrics(ctx context.Context) {
	TycheMetrics = apollo_metrics.InitMetrics()
	TradeMetrics = metrics_trading.NewTradingMetrics(TycheMetrics)
}
