package apollo_ethereum_alerts

import (
	"context"
	"fmt"
	"strings"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/rs/zerolog/log"
)

const hydraLatencyPromQL = "histogram_quantile(0.99, sum(rate(hydra_request_duration_seconds_bucket[30m])) by (le, namespace))"

func (a *ApolloEthereumAlerts) HydraLatencyIssueTrigger(ctx context.Context) error {
	timeNow := time.Now().UTC()
	query := hydraLatencyPromQL
	opts := v1.WithTimeout(time.Second * 10)
	r, w, err := a.Query(ctx, query, timeNow, opts)
	if w != nil {
		log.Warn().Interface("warnings", w).Msg("apollo_ethereum_alerts.BeaconSyncIssueTrigger")
	}
	if err != nil {
		return err
	}
	mv := r.(model.Vector)

	var issuesSlice []string
	for _, v := range mv {
		ns, ok := v.Metric["namespace"]
		if ok && v.Value > 0.5 {
			issuesSlice = append(issuesSlice, fmt.Sprintf("99th percentile latency exceeding %fs in namespace: %s", 0.5, ns))
		}
	}

	if len(issuesSlice) == 0 {
		return nil
	}
	issuesString := strings.Join(issuesSlice, ",")
	alertResp, alertErr := a.CreateLowPriorityAlertFromTemplate(ctx, "This is hydra latency issue alert", "hydra",
		issuesString)
	if alertErr != nil {
		log.Ctx(ctx).Error().Err(alertErr).Msg("apollo_ethereum_alerts.HydraLatencyIssueTrigger")
		return alertErr
	}
	log.Info().Interface("alert", alertResp).Msg("apollo_ethereum_alerts.HydraLatencyIssueTrigger")
	return nil
}
