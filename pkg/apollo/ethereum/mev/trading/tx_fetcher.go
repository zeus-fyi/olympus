package metrics_trading

import (
	"github.com/prometheus/client_golang/prometheus"
)

type TxFetcherMetrics struct {
	TradeMethodStats *prometheus.GaugeVec
	//CurrencyStats    *prometheus.GaugeVec
}

func NewTxFetcherMetrics(reg prometheus.Registerer) TxFetcherMetrics {
	tx := TxFetcherMetrics{}
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
