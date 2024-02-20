package ai_platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
)

func (z *ZeusAiPlatformServicesWorker) ExecuteJsonOutputTaskWorkflow(ctx context.Context, mb *MbChildSubProcessParams) (*MbChildSubProcessParams, error) {
	tc := z.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: z.TaskQueueName,
		ID:        fmt.Sprintf("json-output-task-%s", uuid.New().String()),
	}
	txWf := NewZeusPlatformServiceWorkflows()
	wf := txWf.JsonOutputTaskWorkflow
	workflowRun, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, mb)
	if err != nil {
		log.Err(err).Msg("ExecuteJsonOutputTaskWorkflow")
		return nil, err
	}
	err = workflowRun.Get(ctx, &mb)
	if err != nil {
		log.Err(err).Msg("ExecuteJsonOutputTaskWorkflow: Get ChatCompletionQueryResponse")
		return nil, err
	}
	return mb, nil
}
