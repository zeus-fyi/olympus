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

func (k *KronosActivities) GetAssignments(ctx context.Context) (artemis_orchestrations.OrchestrationJob, error) {
	oj := artemis_orchestrations.OrchestrationJob{}
	err := oj.SelectOrchestrationsWithInstructions(context.Background())
	if err != nil {
		return oj, err
	}
	return oj, err
}
