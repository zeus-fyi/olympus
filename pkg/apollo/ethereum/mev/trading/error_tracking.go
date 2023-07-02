package metrics_trading

import "github.com/prometheus/client_golang/prometheus"

type ErrTrackingMetrics struct {
	PricingData *prometheus.GaugeVec
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
	reg.MustRegister(tx.PricingData)
	return tx
}

func (t *ErrTrackingMetrics) RecordError(method, pair string) {
	t.PricingData.WithLabelValues(pair, method).Add(1)
}
