package metrics_trading

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

type TradeAnalysisMetrics struct {
	SandwichCalculatedRevenueGauge         *prometheus.GaugeVec
	SandwichCalculatedRevenueHistogram     *prometheus.HistogramVec
	SandwichCalculatedROIHistogram         *prometheus.HistogramVec
	SandwichCalculatedUpFrontCostHistogram *prometheus.HistogramVec
}

func NewTradeAnalysisMetrics(reg prometheus.Registerer) TradeAnalysisMetrics {
	tx := TradeAnalysisMetrics{}
	tx.SandwichCalculatedRevenueGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "eth_mev_sandwich_calculated_revenue",
			Help: "Calculated gas free profit from sandwich attack in USD denominated value",
		},
		[]string{"method", "in"},
	)
	tx.SandwichCalculatedRevenueHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "eth_mev_sandwich_calculated_revenue_event",
			Help: "Calculated gas free profit from sandwich attack in USD denominated value",
		},
		[]string{"method", "in"},
	)
	tx.SandwichCalculatedROIHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "eth_mev_sandwich_upfront_cost_event",
			Help: "Calculated upfront cost to execute sandwich attack in USD denominated value",
		},
		[]string{"method", "in"},
	)
	tx.SandwichCalculatedUpFrontCostHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "eth_mev_sandwich_calculated_roi_event",
			Help: "Return on investment from sandwich attack",
		},
		[]string{"method", "in"},
	)
	reg.MustRegister(tx.SandwichCalculatedRevenueGauge, tx.SandwichCalculatedRevenueHistogram, tx.SandwichCalculatedROIHistogram, tx.SandwichCalculatedUpFrontCostHistogram)
	return tx
}

func (t *TradeAnalysisMetrics) CalculatedSandwich(method, in string, upfrontCost, revenue float64) {
	t.SandwichCalculatedRevenueGauge.WithLabelValues(method, in).Add(revenue)
	t.SandwichCalculatedRevenueHistogram.WithLabelValues(method, in).Observe(revenue)
	t.SandwichCalculatedUpFrontCostHistogram.WithLabelValues(method, in).Observe(upfrontCost)
	if upfrontCost == 0 {
		log.Warn().Msg("upfront cost is 0")
		return
	}
	roi := ((revenue / upfrontCost) - 1) * 100
	t.SandwichCalculatedROIHistogram.WithLabelValues(method, in).Observe(roi)
}
