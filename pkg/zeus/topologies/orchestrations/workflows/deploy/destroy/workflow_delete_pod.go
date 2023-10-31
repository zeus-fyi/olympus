package destroy_deployed_workflow

import (
	"time"

	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (t *DestroyDeployTopologyWorkflow) DeletePodWorkflow(ctx workflow.Context, wfId, podName string, cctx zeus_common_types.CloudCtxNs, delay time.Duration) error {
	logger := workflow.GetLogger(ctx)
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second * 60,
		BackoffCoefficient: 1.0,
		MaximumInterval:    time.Minute * 10,
	}
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
		RetryPolicy:         retryPolicy,
	}
	err := workflow.Sleep(ctx, delay)
	if err != nil {
		return err
	}

	dpCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(dpCtx, t.DestroyDeployTopologyActivities.DeletePod, podName, cctx).Get(dpCtx, nil)
	if err != nil {
		logger.Error("Failed to delete pod", "Error", err)
		return err
	}
	return nil
}
