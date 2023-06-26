package metrics_txfetcher

import (
	"github.com/prometheus/client_golang/prometheus"
)

type TxFetcherMetrics struct {
	UniqueTxsCount prometheus.Gauge
}

func NewTxFetcherMetrics() TxFetcherMetrics {
	tx := TxFetcherMetrics{}
	tx.UniqueTxsCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ethereum_unique_txs_count",
		Help: "Count of unique txs added by the tx fetcher",
	})
	return tx
}
