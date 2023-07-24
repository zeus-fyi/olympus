package metrics_trading

import "github.com/prometheus/client_golang/prometheus"

type ErrTrackingMetrics struct {
	HigherThanMaxTradeSizeErrCount prometheus.Counter
	PricingData                    *prometheus.GaugeVec
}

func NewErrTrackingMetrics(reg prometheus.Registerer) ErrTrackingMetrics {
	tx := ErrTrackingMetrics{}
	tx.PricingData = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "eth_pricing_data_err_stats",
			Help: "Tracks errors from fetching pricing data",
		},
		[]string{"pair", "method"},
	)
	tx.HigherThanMaxTradeSizeErrCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "eth_mev_trade_size_too_high_count",
			Help: "Counts the number of times a trade size is higher than max trade size setting",
		},
	)
	reg.MustRegister(tx.PricingData, tx.HigherThanMaxTradeSizeErrCount)
	return tx
}

func (t *ErrTrackingMetrics) CountTradeSizeErr() {
	t.HigherThanMaxTradeSizeErrCount.Add(1)
}

func (t *ErrTrackingMetrics) RecordError(method, pair string) {
	t.PricingData.WithLabelValues(pair, method).Add(1)
}
