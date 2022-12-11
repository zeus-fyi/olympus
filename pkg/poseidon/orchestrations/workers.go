package orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
)

func (t *PoseidonWorker) ExecutePoseidonSyncWorkflow(ctx context.Context, params interface{}) error {
	c := t.ConnectTemporalClient()
	defer c.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	txWf := NewPoseidonSyncWorkflow()
	wf := txWf.Workflow

	_, err := c.ExecuteWorkflow(ctx, workflowOptions, wf, params)
	if err != nil {
		log.Err(err).Msg("ExecutePoseidonSyncWorkflow")
		return err
	}
	return err
}
