package mev_promql

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/tidwall/pretty"
	apollo_prometheus "github.com/zeus-fyi/olympus/pkg/apollo/prometheus"
)

type MevPromQL struct {
	printOn bool
	pc      apollo_prometheus.Prometheus
}

func NewMevPromQL(p apollo_prometheus.Prometheus) MevPromQL {
	return MevPromQL{false, p}
}

var ProxyMevPromQL MevPromQL

// <service-name>.<namespace>.svc.cluster.local
// http://prometheus-operated.observability.svc.cluster.local:9090
// topk(15,sum(eth_mempool_mev_currency_in_stats) by (in))

func (m *MevPromQL) GetTopTokens(ctx context.Context, window v1.Range) ([]Metrics, error) {
	query := fmt.Sprintf("topk(15,sum(eth_mempool_mev_currency_in_stats) by (in))")
	opts := v1.WithTimeout(time.Second * 10)
	r, _, err := m.pc.QueryRange(ctx, query, window, opts)
	if err != nil {
		return nil, err
	}
	var metrics []Metrics
	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &metrics)
	if err != nil {
		return nil, err
	}
	if m.printOn {
		fmt.Println(r)
		requestJSON := pretty.Pretty(b)
		requestJSON = pretty.Color(requestJSON, pretty.TerminalStyle)
		fmt.Println(string(requestJSON))
	}
	return metrics, err
}

func (m *MevPromQL) GetTopRevenuePairs(ctx context.Context, window v1.Range) ([]Metrics, error) {
	query := fmt.Sprintf("topk(25,sum(eth_mev_sandwich_calculated_revenue_event_sum) by (in, pair))")
	opts := v1.WithTimeout(time.Second * 10)
	r, _, err := m.pc.QueryRange(ctx, query, window, opts)
	if err != nil {
		return nil, err
	}
	var metrics []Metrics
	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &metrics)
	if err != nil {
		return nil, err
	}
	if m.printOn {
		fmt.Println(r)
		requestJSON := pretty.Pretty(b)
		requestJSON = pretty.Color(requestJSON, pretty.TerminalStyle)
		fmt.Println(string(requestJSON))
	}
	return metrics, err
}

type Metrics struct {
	Metric struct {
		In   string `json:"in"`
		Pair string `json:"pair"`
	} `json:"metric"`
	Values []interface{} `json:"values"`
}

type TokenMetrics struct {
	Metric struct {
		In string `json:"in"`
	} `json:"metric"`
	Values []interface{} `json:"values"`
}
