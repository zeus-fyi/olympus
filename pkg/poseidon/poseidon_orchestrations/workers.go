package poseidon_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	pg_poseidon "github.com/zeus-fyi/olympus/datastores/postgres/apps/poseidon"
	"go.temporal.io/sdk/client"
)

func (t *PoseidonWorker) ExecutePoseidonSyncWorkflow(ctx context.Context) error {
	c := t.ConnectTemporalClient()
	defer c.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	// 	CronSchedule: "15 8 * * *",
	psWf := NewPoseidonSyncWorkflow(PoseidonSyncActivitiesOrchestrator)
	wf := psWf.PoseidonEthereumWorkflow
	_, err := c.ExecuteWorkflow(ctx, workflowOptions, wf, nil)
	if err != nil {
		log.Err(err).Msg("ExecutePoseidonSyncWorkflow")
		return err
	}
	return err
}

func (t *PoseidonWorker) ExecutePoseidonDiskWipeWorkflow(ctx context.Context, dw pg_poseidon.DiskWipeOrchestration) error {
	c := t.ConnectTemporalClient()
	defer c.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	psWf := NewPoseidonSyncWorkflow(PoseidonSyncActivitiesOrchestrator)
	wf := psWf.PoseidonEthereumClientDiskWipeWorkflow
	_, err := c.ExecuteWorkflow(ctx, workflowOptions, wf, dw)
	if err != nil {
		log.Err(err).Msg("ExecutePoseidonDiskWipeWorkflow")
		return err
	}
	return err
}

func (t *PoseidonWorker) ExecutePoseidonDiskUploadWorkflow(ctx context.Context, du pg_poseidon.UploadDataDirOrchestration) error {
	c := t.ConnectTemporalClient()
	defer c.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	psWf := NewPoseidonSyncWorkflow(PoseidonSyncActivitiesOrchestrator)
	wf := psWf.PoseidonEthereumClientDiskUploadWorkflow
	_, err := c.ExecuteWorkflow(ctx, workflowOptions, wf, du)
	if err != nil {
		log.Err(err).Msg("ExecutePoseidonDiskUploadWorkflow")
		return err
	}
	return err
}

func (t *PoseidonWorker) ExecutePoseidonEthereumClientBeaconUploadWorkflow(ctx context.Context, execClient pg_poseidon.UploadDataDirOrchestration, consensusClient pg_poseidon.UploadDataDirOrchestration) error {
	c := t.ConnectTemporalClient()
	defer c.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	psWf := NewPoseidonSyncWorkflow(PoseidonSyncActivitiesOrchestrator)
	wf := psWf.PoseidonEthereumClientBeaconUploadWorkflow
	_, err := c.ExecuteWorkflow(ctx, workflowOptions, wf, execClient, consensusClient)
	if err != nil {
		log.Err(err).Msg("ExecutePoseidonEthereumClientBeaconUploadWorkflow")
		return err
	}
	return err
}

/*
	// CronSchedule - Optional cron schedule for workflow. If a cron schedule is specified, the workflow will run
	// as a cron based on the schedule. The scheduling will be based on UTC time. Schedule for next run only happen
	// after the current run is completed/failed/timeout. If a RetryPolicy is also supplied, and the workflow failed
	// or timeout, the workflow will be retried based on the retry policy. While the workflow is retrying, it won't
	// schedule its next run. If next schedule is due while workflow is running (or retrying), then it will skip that
	// schedule. Cron workflow will not stop until it is terminated or canceled (by returning temporal.CanceledError).
	// The cron spec is as following:
	// ┌───────────── minute (0 - 59)
	// │ ┌───────────── hour (0 - 23)
	// │ │ ┌───────────── day of the month (1 - 31)
	// │ │ │ ┌───────────── month (1 - 12)
	// │ │ │ │ ┌───────────── day of the week (0 - 6) (Sunday to Saturday)
	// │ │ │ │ │
	// │ │ │ │ │
	// * * * * *
	CronSchedule string
*/
