package metrics_trading

import "github.com/prometheus/client_golang/prometheus"

type TradingMetrics struct {
	TxFetcherMetrics TxFetcherMetrics
}

func NewTradingMetrics(reg prometheus.Registerer) TradingMetrics {
	return TradingMetrics{
		TxFetcherMetrics: NewTxFetcherMetrics(reg),
	}
}
