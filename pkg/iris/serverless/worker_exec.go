package iris_serverless

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	"go.temporal.io/sdk/client"
)

func (i *IrisServicesWorker) ExecuteIrisServerlessResyncWorkflow(ctx context.Context) error {
	tc := i.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		ID:        uuid.New().String(),
		TaskQueue: i.TaskQueueName,
	}
	txWf := NewIrisPlatformServiceWorkflows()
	wf := txWf.IrisServerlessResyncWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID)
	if err != nil {
		log.Err(err).Msg("ExecuteIrisServerlessResyncWorkflow")
		return err
	}
	return err
}

func (i *IrisServicesWorker) ExecuteIrisServerlessPodRestartWorkflow(ctx context.Context, cctx zeus_common_types.CloudCtxNs, podName string, delay time.Duration) error {
	tc := i.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		ID:        uuid.New().String(),
		TaskQueue: i.TaskQueueName,
	}
	txWf := NewIrisPlatformServiceWorkflows()
	wf := txWf.IrisServerlessPodRestartWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID, cctx, podName, delay)
	if err != nil {
		log.Err(err).Msg("IrisServerlessPodRestartWorkflow")
		return err
	}
	return err
}
