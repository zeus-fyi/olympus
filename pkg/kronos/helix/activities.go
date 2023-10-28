package kronos_helix

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
)

const (
	internalOrgID = 7138983863666903883
	olympus       = "olympus"
)

type KronosActivities struct{}

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
		k.CheckEndpointHealth,
	}
}

func (k *KronosActivities) Recycle(ctx context.Context) error {
	err := KronosServiceWorker.ExecuteKronosWorkflow(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (k *KronosActivities) UpsertAssignment(ctx context.Context, oj artemis_orchestrations.OrchestrationJob) error {
	err := oj.UpsertOrchestrationWithInstructions(ctx)
	if err != nil {
		log.Err(err).Interface("oj", oj).Msg("UpsertAssignment: UpsertOrchestrationWithInstructions failed")
		return err
	}
	return nil
}

func (k *KronosActivities) GetInternalAssignments(ctx context.Context) ([]artemis_orchestrations.OrchestrationJob, error) {
	ojs, err := artemis_orchestrations.SelectSystemOrchestrationsWithInstructionsByGroup(ctx, internalOrgID, olympus)
	if err != nil {
		log.Err(err).Msg("GetInternalAssignments: SelectSystemOrchestrationsWithInstructionsByGroup failed")
		return nil, err
	}
	return ojs, err
}

func (k *KronosActivities) GetInstructionsFromJob(ctx context.Context, oj artemis_orchestrations.OrchestrationJob) (Instructions, error) {
	ins := Instructions{}
	err := json.Unmarshal([]byte(oj.Instructions), &ins)
	if err != nil {
		log.Err(err).Interface("oj", oj).Msg("GetInstructionsFromJob: json.Unmarshal failed")
		return ins, err
	}
	return ins, nil
}

func (k *KronosActivities) ProcessAssignment(ctx context.Context, oj artemis_orchestrations.OrchestrationJob) error {
	return nil
}

func (k *KronosActivities) UpdateAndMarkOrchestrationInactive(ctx context.Context, oj *artemis_orchestrations.OrchestrationJob) error {
	oj.Active = false
	err := oj.UpdateOrchestrationActiveStatus(ctx)
	if err != nil {
		log.Err(err).Msg("UpdateAndMarkOrchestrationInactive: UpdateAndMarkOrchestrationInactive failed")
		return err
	}
	return err
}

func (k *KronosActivities) UpdateAndMarkOrchestrationActive(ctx context.Context, oj *artemis_orchestrations.OrchestrationJob) error {
	oj.Active = true
	err := oj.UpdateOrchestrationActiveStatus(ctx)
	if err != nil {
		log.Err(err).Msg("UpdateAndMarkOrchestrationActive: UpdateOrchestrationActiveStatus failed")
		return err
	}
	return err
}
