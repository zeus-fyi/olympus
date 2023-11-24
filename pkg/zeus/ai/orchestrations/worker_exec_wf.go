package ai_platform_service_orchestrations

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	hera_openai_dbmodels "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hera_discord "github.com/zeus-fyi/olympus/pkg/hera/discord"
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

func (h *ZeusAiPlatformServicesWorker) ExecuteAiIngestDiscordWorkflow(ctx context.Context, ou org_users.OrgUser, searchGroupName string, cm hera_discord.ChannelMessages) error {
	tc := h.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: h.TaskQueueName,
		ID:        uuid.New().String(),
	}
	txWf := NewZeusPlatformServiceWorkflows()
	wf := txWf.AiIngestDiscordWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID, ou, searchGroupName, cm)
	if err != nil {
		log.Err(err).Msg("ExecuteAiDiscordWorkflow")
		return err
	}
	return nil
}

func (h *ZeusAiPlatformServicesWorker) ExecuteAiFetchDataToIngestDiscordWorkflow(ctx context.Context, ou org_users.OrgUser, searchGroupName string) error {
	tc := h.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: h.TaskQueueName,
		ID:        uuid.New().String(),
	}
	txWf := NewZeusPlatformServiceWorkflows()
	wf := txWf.AiFetchDataToIngestDiscordWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID, ou, searchGroupName)
	if err != nil {
		log.Err(err).Msg("ExecuteAiFetchDataToIngestDiscordWorkflow")
		return err
	}
	return nil
}
