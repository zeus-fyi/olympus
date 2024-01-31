package kronos_helix

import (
	"context"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
)

const (
	TemporalAlerts = "TEMPORAL_ALERTS"
)

func (k *KronosActivities) GetAlertAssignmentFromInstructions(ctx context.Context, ins Instructions) (*pagerduty.V2Event, error) {
	ojs, err := artemis_orchestrations.SelectActiveOrchestrationsWithInstructionsUsingTimeWindow(ctx, internalOrgID, ins.Type, ins.GroupName, ins.Trigger.AlertAfterTime)
	if err != nil {
		log.Err(err).Interface("ojs", ojs).Interface("inst", ins).Msg("GetAlertAssignmentFromInstructions failed")
		return nil, err
	}
	if len(ojs) == 0 {
		return nil, nil
	}
	pdEvent := PdAlertGenericWfIssuesEvent
	pdEvent.DedupKey = uuid.New().String()
	if ins.Alerts.Message != "" {
		pdEvent.Payload.Summary = ins.Alerts.Message
	}
	if ins.Alerts.Component != "" {
		pdEvent.Payload.Component = ins.Alerts.Component
	}
	if ins.Alerts.Source != "" {
		pdEvent.Payload.Details = ins.Alerts.Source
	}
	pdEvent.Payload.Severity = ins.Alerts.Severity.Critical()
	return &pdEvent, err
}

func (k *KronosActivities) ExecuteTriggeredAlert(ctx context.Context, pdEvent *pagerduty.V2Event) error {
	if pdEvent == nil {
		return nil
	}
	if PdAlertClient.Client == nil {
		panic("PdAlertClient is not initialized")
	}
	if pdEvent.RoutingKey == "" {
		panic("pdEvent.RoutingKey is empty")
	}
	resp, err := PdAlertClient.SendAlert(ctx, *pdEvent)
	if err != nil {
		log.Err(err).Interface("resp", resp).Interface("pdEvent", pdEvent).Msg("ExecuteTriggeredAlert: SendAlert failed")
		return err
	}
	return err
}
