package kronos_helix

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
)

type KronosActivities struct {
}

func NewKronosActivities() KronosActivities {
	return KronosActivities{}
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (k *KronosActivities) GetActivities() ActivitiesSlice {
	return []interface{}{k.Recycle, k.GetAssignments, k.ExecuteTriggeredAlert, k.ProcessAssignment}
}

func (k *KronosActivities) Recycle(ctx context.Context) error {
	err := KronosServiceWorker.ExecuteKronosWorkflow(ctx)
	if err != nil {
		return err
	}
	return nil
}

const internalOrgID = 7138983863666903883

func (k *KronosActivities) GetAssignments(ctx context.Context, orchestType string) ([]artemis_orchestrations.OrchestrationJob, error) {
	ojs, err := artemis_orchestrations.SelectActiveOrchestrationsWithInstructions(ctx, internalOrgID, orchestType)
	if err != nil {
		return nil, err
	}
	return ojs, err
}

func (k *KronosActivities) ProcessAssignment(ctx context.Context, oj artemis_orchestrations.OrchestrationJob) error {
	if oj.Instructions == "{}" || oj.Instructions == "" {
		log.Ctx(ctx).Info().Msg("ProcessAssignment: Instructions are empty")
		return nil
	}
	return nil
}

func (k *KronosActivities) ExecuteTriggeredAlert(ctx context.Context, instructions Instructions) error {
	pdEvent := PdAlertGenericWfIssuesEvent
	pdEvent.Payload.Summary = instructions.Alerts.AlertMessage
	pdEvent.Payload.Component = instructions.Alerts.Component
	pdEvent.Payload.Details = instructions.Alerts.Source
	pdEvent.Payload.Severity = instructions.Alerts.AlertSeverity.Critical()
	_, err := PdAlertClient.SendAlert(ctx, pdEvent)
	if err != nil {
		log.Err(err).Msg("ExecuteTriggeredAlert: SendAlert failed")
		return err
	}
	return err
}

// TODO, activity for resolving alert
