package apollo_ethereum_alerts

import (
	"context"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/rs/zerolog/log"
)

const syncTriggerPromQL = "sum(validator_monitor_slashed) > 0"

func (a *ApolloEthereumAlerts) BeaconSyncIssueTrigger(ctx context.Context) (bool, error) {
	timeNow := time.Now().UTC()
	window := v1.Range{
		Start: timeNow.Add(-time.Minute * 60),
		End:   time.Now().UTC(),
		Step:  time.Minute,
	}
	query := syncTriggerPromQL
	opts := v1.WithTimeout(time.Second * 10)
	r, w, err := a.QueryRange(ctx, query, window, opts)
	if w != nil {
		log.Warn().Interface("warnings", w).Msg("apollo_ethereum_alerts.BeaconSyncIssueTrigger")
	}
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("apollo_ethereum_alerts.BeaconSyncIssueTrigger")
		return false, err
	}
	if len(r.String()) <= 0 {
		return false, nil
	}
	alertResp, alertErr := a.CreateAlertFromTemplate(ctx, "This is a critical ethereum beacon sync issue alert", "beacons",
		"One or more beacons have been detected to have a critical sync issue like a corrupted database")
	if alertErr != nil {
		log.Ctx(ctx).Error().Err(alertErr).Msg("apollo_ethereum_alerts.BeaconSyncIssueTrigger")
		return true, alertErr
	}
	log.Info().Interface("alert", alertResp).Msg("apollo_ethereum_alerts.BeaconSyncIssueTrigger")
	return true, alertErr
}
