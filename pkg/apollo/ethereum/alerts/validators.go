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

const slashingTriggerPromQL = "sum by (namespace, validator) (validator_monitor_slashed > 0)"

func (a *ApolloEthereumAlerts) SlashingAlertTrigger(ctx context.Context) (bool, error) {
	timeNow := time.Now().UTC()
	window := v1.Range{
		Start: timeNow.Add(-time.Minute * 60),
		End:   time.Now().UTC(),
		Step:  time.Minute,
	}
	query := slashingTriggerPromQL
	opts := v1.WithTimeout(time.Second * 10)
	r, w, err := a.QueryRange(ctx, query, window, opts)
	if w != nil {
		log.Warn().Interface("warnings", w).Msg("apollo_ethereum_alerts.SlashingAlert")
	}
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("apollo_ethereum_alerts.SlashingAlert")
		return false, err
	}

	mv := r.(model.Matrix)
	if len(mv) <= 0 {
		return false, nil
	}

	metricStrings := make([]string, len(mv))
	for i, met := range mv {
		metricStrings[i] = met.Metric.String()
	}

	issueString := strings.Join(metricStrings, ",")

	alertResp, alertErr := a.CreateAlertFromTemplate(ctx,
		"IMMEDIATE ACTION NEEDED: This is a critical ethereum validators slashing alert",
		"validators",
		fmt.Sprintf("One or more validators have been slashed: %s", issueString))
	if alertErr != nil {
		log.Ctx(ctx).Error().Err(alertErr).Msg("apollo_ethereum_alerts.SlashingAlert")
		return true, alertErr
	}
	log.Info().Interface("alert", alertResp).Msg("apollo_ethereum_alerts.SlashingAlert")
	return true, alertErr
}
