package poseidon_orchestrations

import (
	"time"

	pg_poseidon "github.com/zeus-fyi/olympus/datastores/postgres/apps/poseidon"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const diskWipeTimeout = 15 * time.Minute

func (t *PoseidonSyncWorkflow) PoseidonEthereumClientDiskWipeWorkflow(ctx workflow.Context, params pg_poseidon.DiskWipeOrchestration) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: diskWipeTimeout,
	}
	aoSync := workflow.ActivityOptions{
		StartToCloseTimeout: diskWipeTimeout,
	}
	syncStatusCheckRetryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second * 60,
		BackoffCoefficient: 2,
	}
	aoSync.RetryPolicy = syncStatusCheckRetryPolicy
	setDiskWipeStatusCtx := workflow.WithActivityOptions(ctx, aoSync)
	err := workflow.ExecuteActivity(setDiskWipeStatusCtx, t.ScheduleDiskWipe, params).Get(setDiskWipeStatusCtx, nil)
	if err != nil {
		log.Error("ScheduleDiskWipe: ", err)
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
