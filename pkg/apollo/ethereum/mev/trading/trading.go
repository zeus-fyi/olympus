package metrics_trading

import "github.com/prometheus/client_golang/prometheus"

type TradingMetrics struct {
	ErrTrackingMetrics      ErrTrackingMetrics
	TxFetcherMetrics        TxFetcherMetrics
	TradeAnalysisMetrics    TradeAnalysisMetrics
	StageProgressionMetrics StageProgressionMetrics
}

func NewTradingMetrics(reg prometheus.Registerer) TradingMetrics {
	return TradingMetrics{
		ErrTrackingMetrics:      NewErrTrackingMetrics(reg),
		TxFetcherMetrics:        NewTxFetcherMetrics(reg),
		TradeAnalysisMetrics:    NewTradeAnalysisMetrics(reg),
		StageProgressionMetrics: NewStageProgressionMetrics(reg),
	}
}
