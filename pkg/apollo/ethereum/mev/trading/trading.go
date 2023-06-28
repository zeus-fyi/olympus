package metrics_trading

import "github.com/prometheus/client_golang/prometheus"

type TradingMetrics struct {
	TxFetcherMetrics     TxFetcherMetrics
	TradeAnalysisMetrics TradeAnalysisMetrics
}

func NewTradingMetrics(reg prometheus.Registerer) TradingMetrics {
	return TradingMetrics{
		TxFetcherMetrics:     NewTxFetcherMetrics(reg),
		TradeAnalysisMetrics: NewTradeAnalysisMetrics(reg),
	}
}
