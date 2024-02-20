package ai_platform_service_orchestrations

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
)

func (z *ZeusAiPlatformServicesWorker) ExecuteRetrievalsWorkflow(ctx context.Context, mb *MbChildSubProcessParams) (*MbChildSubProcessParams, error) {
	tc := z.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: z.TaskQueueName,
		ID:        uuid.New().String(),
	}
	mb.WfID = workflowOptions.ID
	txWf := NewZeusPlatformServiceWorkflows()
	wf := txWf.RetrievalsWorkflow
	res, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, mb)
	if err != nil {
		log.Err(err).Msg("ExecuteRetrievalsWorkflow")
		return nil, err
	}
	err = res.Get(ctx, &mb)
	if err != nil {
		log.Err(err).Msg("ExecuteRetrievalsWorkflow")
		return nil, err
	}
	return mb, nil
}
