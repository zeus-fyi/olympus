package metrics_trading

import (
	"github.com/prometheus/client_golang/prometheus"
)

type TxFetcherMetrics struct {
	Stats *prometheus.GaugeVec
}

func NewTxFetcherMetrics(reg prometheus.Registerer) TxFetcherMetrics {
	tx := TxFetcherMetrics{}
	tx.Stats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ethereum_mev_tx_stats",
			Help: "Count of unique mev txs added by the tx fetcher with additional stats",
		},
		[]string{"address", "method"},
	)
	reg.MustRegister(tx.Stats)
	return tx
}

/*
   1. Write loop to metrics
   2. Create trade count metric
       1. WETH, etc base stable currencies
       2. Other tokens
   3. Create tx destination metric
   4. Simtime execution
   5. Create successful tx token address count

{
		Name: "ethereum_unique_txs_count",
		Help: "Count of unique txs added by the tx fetcher",
	})

txfetcher -> sends to tyche
	tyche pipeline:
		1. decode
		2. binary search
		3. simulate
*/
