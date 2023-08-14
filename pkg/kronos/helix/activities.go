package kronos_helix

import (
	"context"

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
	return []interface{}{k.Recycle, k.GetAssignments}
}

func (k *KronosActivities) Recycle(ctx context.Context) error {
	err := KronosServiceWorker.ExecuteKronosWorkflow(ctx)
	if err != nil {
		return err
	}
	return nil
}

const internalOrgID = 7138983863666903883

func (k *KronosActivities) GetAssignments(ctx context.Context) ([]artemis_orchestrations.OrchestrationJob, error) {
	ojs, err := artemis_orchestrations.SelectActiveOrchestrationsWithInstructions(ctx, internalOrgID)
	if err != nil {
		return nil, err
	}
	return ojs, err
}
