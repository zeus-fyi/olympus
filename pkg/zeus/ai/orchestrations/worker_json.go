package ai_platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
)

func (z *ZeusAiPlatformServicesWorker) ExecuteJsonOutputTaskWorkflow(ctx context.Context, tte TaskToExecute) (*ChatCompletionQueryResponse, error) {
	tc := z.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: z.TaskQueueName,
		ID:        fmt.Sprintf("sm-engagement-%s", uuid.New().String()),
	}
	txWf := NewZeusPlatformServiceWorkflows()
	wf := txWf.JsonOutputTaskWorkflow
	var cr *ChatCompletionQueryResponse
	tte.WfID = workflowOptions.ID
	workflowRun, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, tte)
	if err != nil {
		log.Err(err).Msg("ExecuteJsonOutputTaskWorkflow")
		return nil, err
	}
	err = workflowRun.Get(ctx, &cr)
	if err != nil {
		log.Err(err).Msg("ExecuteJsonOutputTaskWorkflow: Get ChatCompletionQueryResponse")
		return nil, err
	}
	return cr, nil
}
