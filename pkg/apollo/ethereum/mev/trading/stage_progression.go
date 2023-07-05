package metrics_trading

import "github.com/prometheus/client_golang/prometheus"

type StageProgressionMetrics struct {
	//PreEntryFilterTxCount    *prometheus.GaugeVec
	//PostEntryFilterTxCount   *prometheus.GaugeVec
	PostProcessFilterTxCount *prometheus.GaugeVec
	PostSimFilterTxCount     *prometheus.GaugeVec
}

func NewStageProgressionMetrics(reg prometheus.Registerer) StageProgressionMetrics {
	tx := StageProgressionMetrics{}
	//tx.PreEntryFilterTxCount = prometheus.NewGaugeVec(
	//	prometheus.GaugeOpts{
	//		Name: "eth_mev_pre_entry_filter_stage_count",
	//		Help: "Tx count before entry filter",
	//	},
	//	[]string{"pair", "in"},
	//)
	//tx.PostEntryFilterTxCount = prometheus.NewGaugeVec(
	//	prometheus.GaugeOpts{
	//		Name: "eth_mev_post_entry_filter_stage_count",
	//		Help: "Tx count post entry filter",
	//	},
	//	[]string{"pair", "in"},
	//)
	tx.PostProcessFilterTxCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "eth_mev_post_process_tx_stage_count",
			Help: "Tx count after process tx stage",
		},
		[]string{"pair", "in"},
	)
	tx.PostSimFilterTxCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "eth_mev_post_sim_tx_filter_stage_count",
			Help: "Tx count before simulation stage",
		},
		[]string{"pair", "in"},
	)
	reg.MustRegister(tx.PostSimFilterTxCount)
	return tx
}

//func (t *StageProgressionMetrics) CountPreEntryFilterTx(pair, in string) {
//	t.PreEntryFilterTxCount.WithLabelValues(pair, in).Add(1)
//}

//func (t *StageProgressionMetrics) CountPostEntryFilterTx(pair, in string) {
//	t.PostEntryFilterTxCount.WithLabelValues(pair, in).Add(1)
//}

func (t *StageProgressionMetrics) CountPostProcessFilterTx(pair, in string) {
	t.PostProcessFilterTxCount.WithLabelValues(pair, in).Add(1)
}

func (t *StageProgressionMetrics) CountPostSimTx(pair, in string) {
	t.PostSimFilterTxCount.WithLabelValues(pair, in).Add(1)
}
