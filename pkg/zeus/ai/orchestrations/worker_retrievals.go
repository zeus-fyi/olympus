package ai_platform_service_orchestrations

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
)

func (z *ZeusAiPlatformServicesWorker) ExecuteRetrievalsWorkflow(ctx context.Context, tte TaskToExecute) error {
	tc := z.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: z.TaskQueueName,
		ID:        uuid.New().String(),
	}
	txWf := NewZeusPlatformServiceWorkflows()
	wf := txWf.RetrievalsWorkflow
	tte.WfID = workflowOptions.ID
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, tte)
	if err != nil {
		log.Err(err).Msg("ExecuteRetrievalsWorkflow")
		return err
	}
	return nil
}
