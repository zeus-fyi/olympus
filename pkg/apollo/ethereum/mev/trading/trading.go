package metrics_trading

import "github.com/prometheus/client_golang/prometheus"

type TradingMetrics struct {
	MevTxStats *prometheus.GaugeVec
}

func NewTradingMetrics() TradingMetrics {
	tx := TradingMetrics{}
	tx.MevTxStats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "",
			Subsystem: "",
			Name:      "ethereum_unique_txs_count",
			Help:      "Count of unique txs added by the tx fetcher with additional stats",
		},
		[]string{"address", "method"},
	)
	return tx
}
