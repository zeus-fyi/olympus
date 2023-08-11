package kronos_helix

import (
	"time"

	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"go.temporal.io/sdk/workflow"
)

type KronosWorkflow struct {
	temporal_base.Workflow
	KronosActivities
}

const defaultTimeout = 6 * time.Second

func NewKronosWorkflow() KronosWorkflow {
	deployWf := KronosWorkflow{
		Workflow:         temporal_base.Workflow{},
		KronosActivities: KronosActivities{},
	}
	return deployWf
}

func (k *KronosActivities) GetWorkflows() []interface{} {
	return []interface{}{}
}

func (k *KronosWorkflow) T(ctx workflow.Context) error {
	return nil
}
