package metrics_trading

import "github.com/prometheus/client_golang/prometheus"

type TradeAnalysisMetrics struct {
	TradeMethodStats *prometheus.GaugeVec
}

func NewTradeAnalysisMetrics(reg prometheus.Registerer) TradeAnalysisMetrics {
	tx := TradeAnalysisMetrics{}
	tx.TradeMethodStats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "eth_mempool_mev_tx_stats",
			Help: "Count of unique incoming mev txs added by the tx fetcher with additional stats",
		},
		[]string{"address", "method"},
	)
	reg.MustRegister(tx.TradeMethodStats)
	return tx
}

/*
	todo, sandwich profit stats
*/
