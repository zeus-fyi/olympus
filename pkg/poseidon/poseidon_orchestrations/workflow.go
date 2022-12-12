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

const defaultTimeout = 12 * time.Hour

func NewPoseidonSyncWorkflow(psa PoseidonSyncActivities) PoseidonSyncWorkflow {
	deployWf := PoseidonSyncWorkflow{
		Workflow:               temporal_base.Workflow{},
		PoseidonSyncActivities: psa,
	}
	return deployWf
}

func (t *PoseidonSyncWorkflow) GetWorkflows() []interface{} {
	return []interface{}{t.PoseidonEthereumWorkflow}
}

func (t *PoseidonSyncWorkflow) PoseidonEthereumWorkflow(ctx workflow.Context, params interface{}) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	aoSync := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	syncStatusCheckRetryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second * 15,
		BackoffCoefficient: 2,
	}
	aoSync.RetryPolicy = syncStatusCheckRetryPolicy
	execSyncStatusCtx := workflow.WithActivityOptions(ctx, aoSync)
	err := workflow.ExecuteActivity(execSyncStatusCtx, t.IsExecClientSynced).Get(execSyncStatusCtx, nil)
	if err != nil {
		log.Error("IsExecClientSynced: ", err)
		return err
	}
	pauseExecClientCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(pauseExecClientCtx, t.PauseExecClient).Get(pauseExecClientCtx, nil)
	if err != nil {
		log.Error("PauseExecClient: ", err)
		return err
	}
	err = workflow.Sleep(pauseExecClientCtx, 2*time.Minute)
	if err != nil {
		log.Error("PauseExecClient: ", err)
		return err
	}
	rsyncExecClientCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(rsyncExecClientCtx, t.RsyncExecBucket).Get(rsyncExecClientCtx, nil)
	if err != nil {
		log.Error("RsyncExecBucket: ", err)
		return err
	}
	resumeExecClientCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(resumeExecClientCtx, t.ResumeExecClient).Get(resumeExecClientCtx, nil)
	if err != nil {
		log.Error("ResumeExecClient: ", err)
		return err
	}
	consensusSyncStatusCtx := workflow.WithActivityOptions(ctx, aoSync)
	err = workflow.ExecuteActivity(consensusSyncStatusCtx, t.IsConsensusClientSynced).Get(consensusSyncStatusCtx, nil)
	if err != nil {
		log.Error("IsConsensusClientSynced: ", err)
		return err
	}
	pauseConsensusClientCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(pauseConsensusClientCtx, t.PauseConsensusClient).Get(pauseConsensusClientCtx, nil)
	if err != nil {
		log.Error("PauseConsensusClient: ", err)
		return err
	}
	err = workflow.Sleep(pauseConsensusClientCtx, 2*time.Minute)
	if err != nil {
		log.Error("PauseConsensusClient: ", err)
		return err
	}
	rsyncConsensusClientCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(rsyncConsensusClientCtx, t.RsyncConsensusBucket).Get(rsyncConsensusClientCtx, nil)
	if err != nil {
		log.Error("RsyncConsensusBucket: ", err)
		return err
	}
	resumeConsensusClientCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(resumeConsensusClientCtx, t.ResumeConsensusClient).Get(resumeConsensusClientCtx, nil)
	if err != nil {
		log.Error("RsyncConsensusBucket: ", err)
		return err
	}
	return nil
}
