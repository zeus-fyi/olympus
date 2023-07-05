package metrics_trading

import "github.com/prometheus/client_golang/prometheus"

type StageProgressionMetrics struct {
	PreEntryFilterTxCount    prometheus.Counter
	PostEntryFilterTxCount   prometheus.Counter
	PostDecodeTxCount        prometheus.Counter
	PostProcessFilterTxCount prometheus.Counter
	PostSimFilterTxCount     prometheus.Counter
	SavedMempoolTxCount      prometheus.Counter
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
	tx.PostProcessFilterTxCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "eth_mev_post_process_tx_stage_count",
			Help: "Tx count post process tx",
		},
	)
	tx.PostSimFilterTxCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "eth_mev_post_sim_tx_filter_stage",
			Help: "Tx count before simulation stage",
		},
	)
	tx.SavedMempoolTxCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "eth_mev_mempool_tx_recorded",
			Help: "Tx count of mempool tx meeting criteria",
		},
	)
	reg.MustRegister(tx.PostSimFilterTxCount, tx.PostProcessFilterTxCount, tx.PostDecodeTxCount, tx.PostEntryFilterTxCount, tx.PreEntryFilterTxCount, tx.SavedMempoolTxCount)
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

func (t *StageProgressionMetrics) CountPostProcessTx(count float64) {
	t.PostProcessFilterTxCount.Add(count)
}

func (t *StageProgressionMetrics) CountPostSimFilterTx(count float64) {
	t.PostSimFilterTxCount.Add(count)
}

func (t *StageProgressionMetrics) CountSavedMempoolTx(count float64) {
	t.SavedMempoolTxCount.Add(count)
}
