package kronos_helix

import (
	"context"
	"encoding/json"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/google/uuid"
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
	return []interface{}{
		k.Recycle,
		k.GetInternalAssignments,
		k.GetAlertAssignmentFromInstructions,
		k.ExecuteTriggeredAlert,
		k.ProcessAssignment,
		k.UpsertAssignment,
		k.UpdateAndMarkOrchestrationInactive,
		k.UpdateAndMarkOrchestrationActive,
		k.GetInstructionsFromJob,
	}
}

func (k *KronosActivities) Recycle(ctx context.Context) error {
	err := KronosServiceWorker.ExecuteKronosWorkflow(ctx)
	if err != nil {
		return err
	}
	return nil
}

const (
	internalOrgID = 7138983863666903883
	olympus       = "olympus"
)

func (k *KronosActivities) UpsertAssignment(ctx context.Context, oj artemis_orchestrations.OrchestrationJob) error {
	err := oj.UpsertOrchestrationWithInstructions(ctx)
	if err != nil {
		log.Err(err).Msg("UpsertAssignment: UpsertOrchestrationWithInstructions failed")
		return err
	}
	return nil
}

func (k *KronosActivities) GetInternalAssignments(ctx context.Context) ([]artemis_orchestrations.OrchestrationJob, error) {
	ojs, err := artemis_orchestrations.SelectSystemOrchestrationsWithInstructionsByGroup(ctx, internalOrgID, olympus)
	if err != nil {
		return nil, err
	}
	return ojs, err
}

func (k *KronosActivities) GetInstructionsFromJob(ctx context.Context, oj artemis_orchestrations.OrchestrationJob) (Instructions, error) {
	ins := Instructions{}
	err := json.Unmarshal([]byte(oj.Instructions), &ins)
	if err != nil {
		return ins, err
	}
	return ins, nil
}

func (k *KronosActivities) GetAlertAssignmentFromInstructions(ctx context.Context, ins Instructions) (*pagerduty.V2Event, error) {
	ojs, err := artemis_orchestrations.SelectActiveOrchestrationsWithInstructionsUsingTimeWindow(ctx, internalOrgID, ins.Type, ins.GroupName, ins.Trigger.AlertAfterTime)
	if err != nil {
		return nil, err
	}
	if len(ojs) == 0 {
		return nil, nil
	}
	pdEvent := PdAlertGenericWfIssuesEvent
	pdEvent.DedupKey = uuid.New().String()
	pdEvent.Payload.Summary = ins.Alerts.Message
	pdEvent.Payload.Component = ins.Alerts.Component
	pdEvent.Payload.Details = ins.Alerts.Source
	pdEvent.Payload.Severity = ins.Alerts.Severity.Critical()
	return &pdEvent, err
}

func (k *KronosActivities) ProcessAssignment(ctx context.Context, oj artemis_orchestrations.OrchestrationJob) error {
	if oj.Instructions == "{}" || oj.Instructions == "" {
		log.Info().Msg("ProcessAssignment: Instructions are empty")
		return nil
	}
	var instructions Instructions
	err := json.Unmarshal([]byte(oj.Instructions), &instructions)
	if err != nil {
		log.Err(err).Msg("ProcessAssignment: Unmarshal failed")
		return err
	}
	return nil
}

func (k *KronosActivities) UpdateAndMarkOrchestrationInactive(ctx context.Context, oj artemis_orchestrations.OrchestrationJob) error {
	oj.Active = false
	err := oj.UpdateOrchestrationActiveStatus(ctx)
	if err != nil {
		log.Err(err).Msg("UpdateAndMarkOrchestrationInactive: UpdateAndMarkOrchestrationInactive failed")
		return err
	}
	return err
}

func (k *KronosActivities) UpdateAndMarkOrchestrationActive(ctx context.Context, oj artemis_orchestrations.OrchestrationJob) error {
	oj.Active = true
	err := oj.UpdateOrchestrationActiveStatus(ctx)
	if err != nil {
		log.Err(err).Msg("UpdateAndMarkOrchestrationActive: UpdateOrchestrationActiveStatus failed")
		return err
	}
	return err
}

func (k *KronosActivities) ExecuteTriggeredAlert(ctx context.Context, pdEvent pagerduty.V2Event) error {
	_, err := PdAlertClient.SendAlert(ctx, pdEvent)
	if err != nil {
		log.Err(err).Msg("ExecuteTriggeredAlert: SendAlert failed")
		return err
	}
	return err
}
