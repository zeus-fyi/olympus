package metrics_trading

import "github.com/prometheus/client_golang/prometheus"

type StageProgressionMetrics struct {
	PreSimTxCount *prometheus.GaugeVec
}

func NewStageProgressionMetrics(reg prometheus.Registerer) StageProgressionMetrics {
	tx := StageProgressionMetrics{}
	tx.PreSimTxCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "eth_mev_pre_sim_tx_stage_count",
			Help: "Tx count before simulation stage",
		},
		[]string{"pair", "in"},
	)
	reg.MustRegister(tx.PreSimTxCount)
	return tx
}

func (t *StageProgressionMetrics) CountPreSimTx(pair, in string) {
	t.PreSimTxCount.WithLabelValues(pair, in).Add(1)
}
