package ai_platform_service_orchestrations

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"go.temporal.io/sdk/client"
)

func (h *ZeusAiPlatformServicesWorker) ExecuteAiTaskWorkflow(ctx context.Context, ou org_users.OrgUser, em, content string) error {
	tc := h.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: h.TaskQueueName,
		ID:        uuid.New().String(),
	}
	txWf := NewZeusPlatformServiceWorkflows()
	wf := txWf.AiWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID, ou, em, content)
	if err != nil {
		log.Err(err).Msg("ExecuteAiTaskWorkflow")
		return err
	}
	return nil
}
