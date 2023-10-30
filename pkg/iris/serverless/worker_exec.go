package iris_serverless

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
)

func (i *IrisServicesWorker) ExecuteIrisServerlessResyncWorkflow(ctx context.Context) error {
	tc := i.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		ID:        uuid.New().String(),
		TaskQueue: i.TaskQueueName,
	}
	txWf := NewHestiaPlatformServiceWorkflows()
	wf := txWf.IrisServerlessResyncWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID)
	if err != nil {
		log.Err(err).Msg("ExecuteIrisServerlessResyncWorkflow")
		return err
	}
	return err
}
