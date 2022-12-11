package poseidon_orchestrations

import (
	"time"

	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type PoseidonSyncWorkflow struct {
	PoseidonSyncActivities
	temporal_base.Workflow
}

const defaultTimeout = 30 * time.Minute

func NewPoseidonSyncWorkflow() PoseidonSyncWorkflow {
	deployWf := PoseidonSyncWorkflow{
		Workflow: temporal_base.Workflow{},
	}
	return deployWf
}

func (t *PoseidonSyncWorkflow) GetWorkflows() []interface{} {
	return []interface{}{t.PoseidonWorkflow}
}

func (t *PoseidonSyncWorkflow) PoseidonWorkflow(ctx workflow.Context) error {
	//log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	sendCtx := workflow.WithActivityOptions(ctx, ao)
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second * 15,
		BackoffCoefficient: 2,
	}
	ao.RetryPolicy = retryPolicy
	rxCtx := workflow.WithActivityOptions(ctx, ao)

	err := workflow.ExecuteActivity(sendCtx, t.SyncExecStatus).Get(sendCtx, nil)
	if err != nil {
		return err
	}
	err = workflow.ExecuteActivity(sendCtx, t.SyncConsensusStatus).Get(sendCtx, nil)
	if err != nil {
		return err
	}
	// TODO make these do consensus + exec via params
	err = workflow.ExecuteActivity(rxCtx, t.Pause).Get(rxCtx, nil)
	if err != nil {
		return err
	}
	err = workflow.ExecuteActivity(rxCtx, t.Resume).Get(rxCtx, nil)
	if err != nil {
		return err
	}
	err = workflow.ExecuteActivity(rxCtx, t.RsyncExecBucket).Get(rxCtx, nil)
	if err != nil {
		return err
	}
	err = workflow.ExecuteActivity(rxCtx, t.RsyncConsensusBucket).Get(rxCtx, nil)
	if err != nil {
		return err
	}

	return nil
}
