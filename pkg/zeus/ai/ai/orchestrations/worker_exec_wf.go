package ai_platform_service_orchestrations

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	hera_openai_dbmodels "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hermes_email_notifications "github.com/zeus-fyi/olympus/pkg/hermes/email"
	"go.temporal.io/sdk/client"
)

func (h *ZeusAiPlatformServicesWorker) ExecuteAiTaskWorkflow(ctx context.Context, ou org_users.OrgUser, msgs []hermes_email_notifications.EmailContents) error {
	tc := h.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: h.TaskQueueName,
		ID:        uuid.New().String(),
	}
	txWf := NewZeusPlatformServiceWorkflows()
	wf := txWf.AiEmailWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID, ou, msgs)
	if err != nil {
		log.Err(err).Msg("ExecuteAiTaskWorkflow")
		return err
	}
	return nil
}
func (h *ZeusAiPlatformServicesWorker) ExecuteAiTelegramWorkflow(ctx context.Context, ou org_users.OrgUser, msgs []hera_openai_dbmodels.TelegramMessage) error {
	tc := h.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: h.TaskQueueName,
		ID:        uuid.New().String(),
	}
	txWf := NewZeusPlatformServiceWorkflows()
	wf := txWf.AiIngestTelegramWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID, ou, msgs)
	if err != nil {
		log.Err(err).Msg("ExecuteAiTelegramWorkflow")
		return err
	}
	return nil
}

func (h *ZeusAiPlatformServicesWorker) ExecuteAiTwitterWorkflow(ctx context.Context, ou org_users.OrgUser, searchGroupName string) error {
	tc := h.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: h.TaskQueueName,
		ID:        uuid.New().String(),
	}
	txWf := NewZeusPlatformServiceWorkflows()
	wf := txWf.AiIngestTwitterWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID, ou, searchGroupName)
	if err != nil {
		log.Err(err).Msg("AiIngestTwitterWorkflow")
		return err
	}
	return nil
}
func (h *ZeusAiPlatformServicesWorker) ExecuteAiRedditWorkflow(ctx context.Context, ou org_users.OrgUser, searchGroupName string) error {
	tc := h.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: h.TaskQueueName,
		ID:        uuid.New().String(),
	}
	txWf := NewZeusPlatformServiceWorkflows()
	wf := txWf.AiIngestRedditWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID, ou, searchGroupName)
	if err != nil {
		log.Err(err).Msg("ExecuteAiRedditWorkflow")
		return err
	}
	return nil
}
