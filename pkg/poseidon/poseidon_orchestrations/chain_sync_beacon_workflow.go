package poseidon_orchestrations

import (
	"time"

	pg_poseidon "github.com/zeus-fyi/olympus/datastores/postgres/apps/poseidon"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (t *PoseidonSyncWorkflow) PoseidonEthereumClientBeaconUploadWorkflow(ctx workflow.Context,
	execClientParams pg_poseidon.UploadDataDirOrchestration,
	consensusClientParams pg_poseidon.UploadDataDirOrchestration,
) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	aoSync := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	syncStatusCheckRetryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second * 60,
		BackoffCoefficient: 2,
	}
	aoSync.RetryPolicy = syncStatusCheckRetryPolicy
	execSyncStatusCtx := workflow.WithActivityOptions(ctx, aoSync)
	err := workflow.ExecuteActivity(execSyncStatusCtx, t.IsExecClientSynced).Get(execSyncStatusCtx, nil)
	if err != nil {
		log.Error("IsExecClientSynced: ", err)
		return err
	}
	setDiskUploadStatusCtx := workflow.WithActivityOptions(ctx, aoSync)
	err = workflow.ExecuteActivity(setDiskUploadStatusCtx, t.ScheduleDiskUpload, execClientParams).Get(setDiskUploadStatusCtx, nil)
	if err != nil {
		log.Error("ScheduleDiskUpload: ", err)
		return err
	}
	restartCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(restartCtx, t.RestartBeaconPod, execClientParams.ClientName, execClientParams.CloudCtxNs).Get(restartCtx, nil)
	if err != nil {
		log.Error("RestartBeaconPod: ", err)
		return err
	}
	consensusClientSyncStatusCtx := workflow.WithActivityOptions(ctx, aoSync)
	err = workflow.ExecuteActivity(consensusClientSyncStatusCtx, t.IsConsensusClientSynced).Get(consensusClientSyncStatusCtx, nil)
	if err != nil {
		log.Error("IsConsensusClientSynced: ", err)
		return err
	}
	setDiskUploadStatusCtx = workflow.WithActivityOptions(ctx, aoSync)
	err = workflow.ExecuteActivity(setDiskUploadStatusCtx, t.ScheduleDiskUpload, consensusClientParams).Get(setDiskUploadStatusCtx, nil)
	if err != nil {
		log.Error("ScheduleDiskUpload: ", err)
		return err
	}
	restartCtx = workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(restartCtx, t.RestartBeaconPod, consensusClientParams.ClientName, consensusClientParams.CloudCtxNs).Get(restartCtx, nil)
	if err != nil {
		log.Error("RestartBeaconPod: ", err)
		return err
	}
	return nil
}
