package kronos_helix

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
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

func (k *KronosWorkflow) GetWorkflows() []interface{} {
	return []interface{}{k.Yin, k.Yang, k.SignalFlow}
}

// SignalFlow should be used to place new control flows on the helix
func (k *KronosWorkflow) SignalFlow(ctx workflow.Context) error {
	return nil
}

// Yin should send commands and execute actions
func (k *KronosWorkflow) Yin(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
	}
	aCtx := workflow.WithActivityOptions(ctx, ao)
	oj := artemis_orchestrations.OrchestrationJob{}
	err := workflow.ExecuteActivity(aCtx, k.GetAssignments).Get(aCtx, &oj)
	if err != nil {
		return err
	}
	// TODO: Add logic to handle the assignments

	return nil
}

// Yang should check, status, & react to changes
func (k *KronosWorkflow) Yang(ctx workflow.Context) error {
	//workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
	}
	aCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(aCtx, k.Recycle).Get(aCtx, nil)
	if err != nil {
		return err
	}
	return nil
}
