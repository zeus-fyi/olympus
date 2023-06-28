package metrics_trading

import (
	"github.com/prometheus/client_golang/prometheus"
)

type TxFetcherMetrics struct {
	TradeMethodStats *prometheus.GaugeVec
	CurrencyStatsIn  *prometheus.GaugeVec
	CurrencyStatsOut *prometheus.GaugeVec
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
	tx.CurrencyStatsIn = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "eth_mempool_mev_currency_in_stats",
			Help: "Trade currency in and value",
		},
		[]string{"address", "method", "in"},
	)
	tx.CurrencyStatsOut = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "eth_mempool_mev_currency_out_stats",
			Help: "Trade currency out address & value",
		},
		[]string{"address", "method", "out"},
	)
	reg.MustRegister(tx.TradeMethodStats, tx.CurrencyStatsIn, tx.CurrencyStatsOut)
	return tx
}
