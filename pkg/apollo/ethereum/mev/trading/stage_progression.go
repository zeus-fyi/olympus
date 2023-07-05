package metrics_trading

import "github.com/prometheus/client_golang/prometheus"

type StageProgressionMetrics struct {
	PreEntryFilterTxCount    prometheus.Counter
	PostEntryFilterTxCount   prometheus.Counter
	PostDecodeTxCount        prometheus.Counter
	PostProcessFilterTxCount *prometheus.GaugeVec
	PostSimFilterTxCount     *prometheus.GaugeVec
}

func NewStageProgressionMetrics(reg prometheus.Registerer) StageProgressionMetrics {
	tx := StageProgressionMetrics{}
	tx.PreEntryFilterTxCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "eth_mev_pre_entry_filter_stage_count",
			Help: "Tx count before entry filter",
		},
	)
	tx.PostEntryFilterTxCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "eth_mev_post_entry_filter_stage_count",
			Help: "Tx count post entry filter",
		},
	)
	tx.PostDecodeTxCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "eth_mev_post_decode_tx_stage_count",
			Help: "Tx count pre process tx entry",
		},
	)
	tx.PostProcessFilterTxCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "eth_mev_post_process_tx_stage",
			Help: "Tx count after process tx stage",
		},
		[]string{"pair", "in"},
	)
	tx.PostSimFilterTxCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "eth_mev_post_sim_tx_filter_stage",
			Help: "Tx count before simulation stage",
		},
		[]string{"pair", "in"},
	)
	reg.MustRegister(tx.PostSimFilterTxCount)
	return tx
}

func (t *StageProgressionMetrics) CountPostDecodeTx() {
	t.PostDecodeTxCount.Add(1)
}

func (t *StageProgressionMetrics) CountPreEntryFilterTx() {
	t.PreEntryFilterTxCount.Add(1)
}

func (t *StageProgressionMetrics) CountPostEntryFilterTx() {
	t.PostEntryFilterTxCount.Add(1)
}

func (t *StageProgressionMetrics) CountPostProcessFilterTx(pair, in string) {
	t.PostProcessFilterTxCount.WithLabelValues(pair, in).Add(1)
}

func (t *StageProgressionMetrics) CountPostSimTx(pair, in string) {
	t.PostSimFilterTxCount.WithLabelValues(pair, in).Add(1)
}
