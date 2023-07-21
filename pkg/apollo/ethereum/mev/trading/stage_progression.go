package metrics_trading

import "github.com/prometheus/client_golang/prometheus"

type StageProgressionMetrics struct {
	PreEntryFilterTxCount              prometheus.Counter
	PostEntryFilterTxCount             prometheus.Counter
	PostDecodeTxCount                  prometheus.Counter
	PostProcessFilterTxCount           prometheus.Counter
	PostSimFilterTxCount               prometheus.Counter
	PostActiveTradingSimFilterTxCount  prometheus.Counter
	PostSimStageCount                  prometheus.Counter
	SentFlashbotsBundleSubmissionCount prometheus.Counter
	SavedMempoolTxCount                prometheus.Counter

	CheckpointOneMarker prometheus.Counter
	CheckpointTwoMarker prometheus.Counter
}

func NewStageProgressionMetrics(reg prometheus.Registerer) StageProgressionMetrics {
	tx := StageProgressionMetrics{}
	tx.CheckpointOneMarker = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "eth_mev_checkpoint_one_marker",
			Help: "Checkpoint marker for debugging",
		},
	)
	tx.CheckpointTwoMarker = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "eth_mev_checkpoint_two_marker",
			Help: "Checkpoint marker for debugging",
		},
	)
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
	tx.PostActiveTradingSimFilterTxCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "eth_mev_post_active_trading_sim_tx_filter_stage",
			Help: "Tx count before active simulation stage",
		},
	)
	tx.SentFlashbotsBundleSubmissionCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "eth_mev_sent_flashbots_bundle_submission_count",
			Help: "Tx count of flashbots bundle submissions",
		},
	)
	tx.PostSimStageCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "eth_mev_post_sim_stage_count",
			Help: "Tx count after simulation stage",
		},
	)
	reg.MustRegister(
		tx.PostSimFilterTxCount,
		tx.PostProcessFilterTxCount,
		tx.PostDecodeTxCount,
		tx.PostEntryFilterTxCount,
		tx.PreEntryFilterTxCount,
		tx.SavedMempoolTxCount,
		tx.PostSimStageCount,
		tx.PostActiveTradingSimFilterTxCount,
		tx.SentFlashbotsBundleSubmissionCount,
		tx.CheckpointOneMarker,
		tx.CheckpointTwoMarker,
	)
	return tx
}

func (t *StageProgressionMetrics) CountPostDecodeTx() {
	t.PostDecodeTxCount.Add(1)
}

func (t *StageProgressionMetrics) CountCheckpointOneMarker() {
	t.CheckpointOneMarker.Add(1)
}

func (t *StageProgressionMetrics) CountCheckpointTwoMarker() {
	t.CheckpointTwoMarker.Add(1)
}

func (t *StageProgressionMetrics) CountPreEntryFilterTx() {
	t.PreEntryFilterTxCount.Add(1)
}

func (t *StageProgressionMetrics) CountPostEntryFilterTx() {
	t.PostEntryFilterTxCount.Add(1)
}

func (t *StageProgressionMetrics) CountPostActiveTradingFilter(count float64) {
	t.PostActiveTradingSimFilterTxCount.Add(count)
}

func (t *StageProgressionMetrics) CountSentFlashbotsBundleSubmission(count float64) {
	t.SentFlashbotsBundleSubmissionCount.Add(count)
}

func (t *StageProgressionMetrics) CountPostProcessTx(count float64) {
	t.PostProcessFilterTxCount.Add(count)
}

func (t *StageProgressionMetrics) CountPostSimFilterTx(count float64) {
	t.PostSimFilterTxCount.Add(count)
}

func (t *StageProgressionMetrics) CountPostSimStage(count float64) {
	t.PostSimStageCount.Add(count)
}

func (t *StageProgressionMetrics) CountSavedMempoolTx(count float64) {
	t.SavedMempoolTxCount.Add(count)
}
