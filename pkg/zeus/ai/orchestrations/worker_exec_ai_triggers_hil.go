package ai_platform_service_orchestrations

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
)

func (z *ZeusAiPlatformServicesWorker) ExecuteTriggerActionsWorkflow(ctx context.Context, approvalTaskGroup ApprovalTaskGroup) error {
	tc := z.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: z.TaskQueueName,
		ID:        uuid.New().String(),
	}
	txWf := NewZeusPlatformServiceWorkflows()
	wf := txWf.TriggerActionsWorkflow
	approvalTaskGroup.WfID = workflowOptions.ID
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, approvalTaskGroup)
	if err != nil {
		log.Err(err).Msg("ExecuteTriggerActionsWorkflow")
		return err
	}
	return nil
}
