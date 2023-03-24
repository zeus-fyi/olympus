package poseidon_orchestrations

import (
	"time"

	pg_poseidon "github.com/zeus-fyi/olympus/datastores/postgres/apps/poseidon"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (t *PoseidonSyncWorkflow) PoseidonEthereumClientDiskUploadWorkflow(ctx workflow.Context, params pg_poseidon.UploadDataDirOrchestration) error {
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
	setDiskUploadStatusCtx := workflow.WithActivityOptions(ctx, aoSync)
	err := workflow.ExecuteActivity(setDiskUploadStatusCtx, t.ScheduleUploadDisk, params).Get(setDiskUploadStatusCtx, nil)
	if err != nil {
		log.Error("ScheduleDiskUpload: ", err)
		return err
	}
	restartCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(restartCtx, t.RestartBeaconPod, params).Get(restartCtx, nil)
	if err != nil {
		log.Error("RestartBeaconPod: ", err)
		return err
	}
	return nil
}
