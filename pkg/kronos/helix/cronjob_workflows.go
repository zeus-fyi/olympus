package kronos_helix

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (k *KronosWorkflow) KronosCronJob(ctx workflow.Context, mi Instructions, pollCycles int) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    mi.CronJob.PollInterval,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
		},
	}
	for i := 0; i < pollCycles; i++ {
		cronJobStartCtx := workflow.WithActivityOptions(ctx, ao)
		err := workflow.ExecuteActivity(cronJobStartCtx, k.StartCronJobWorkflow, mi).Get(cronJobStartCtx, nil)
		if err != nil {
			logger.Error("failed to execute cronjob", "Error", err)
			return err
		}
		err = workflow.Sleep(ctx, mi.CronJob.PollInterval)
		if err != nil {
			logger.Error("failed to sleep", "Error", err)
			return err
		}
	}
	return nil
}
