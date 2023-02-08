package apollo_metrics

import "github.com/prometheus/client_golang/prometheus"

func InitMetrics(cs ...prometheus.Collector) *prometheus.Registry {
	registry := prometheus.NewRegistry()
	registry.MustRegister(cs...)
	return registry
}
