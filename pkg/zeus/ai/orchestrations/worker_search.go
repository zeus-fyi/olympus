package ai_platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	hera_openai_dbmodels "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hera_discord "github.com/zeus-fyi/olympus/pkg/hera/discord"
	"go.temporal.io/sdk/client"
)

func (z *ZeusAiPlatformServicesWorker) ExecuteAiSearchIndexerActionsWorkflow(ctx context.Context, ou org_users.OrgUser, params SearchIndexerActionsRequest) error {
	tc := z.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: z.TaskQueueName,
		ID:        fmt.Sprintf("search-indexer-actions-%s", uuid.New().String()),
	}
	txWf := NewZeusPlatformServiceWorkflows()
	wf := txWf.AiSearchIndexerActionsWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID, ou, params)
	if err != nil {
		log.Err(err).Msg("ExecuteAiSearchIndexerActionsWorkflow")
		return err
	}
	return nil
}

func (z *ZeusAiPlatformServicesWorker) ExecuteAiSearchIndexerWorkflow(ctx context.Context) error {
	tc := z.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: z.TaskQueueName,
		ID:        fmt.Sprintf("search-indexer-%s", uuid.New().String()),
	}
	txWf := NewZeusPlatformServiceWorkflows()
	wf := txWf.AiSearchIndexerWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID)
	if err != nil {
		log.Err(err).Msg("ExecuteAiSearchIndexerWorkflow")
		return err
	}
	return nil
}

func (z *ZeusAiPlatformServicesWorker) ExecuteAiTelegramWorkflow(ctx context.Context, ou org_users.OrgUser, msgs []hera_openai_dbmodels.TelegramMessage) error {
	tc := z.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: z.TaskQueueName,
		ID:        fmt.Sprintf("telegram-%s", uuid.New().String()),
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

func (z *ZeusAiPlatformServicesWorker) ExecuteAiTwillioWorkflow(ctx context.Context, ou org_users.OrgUser, msgs []hera_openai_dbmodels.TelegramMessage) error {
	tc := z.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: z.TaskQueueName,
		ID:        fmt.Sprintf("twillio-%s", uuid.New().String()),
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

func (z *ZeusAiPlatformServicesWorker) ExecuteAiTwitterWorkflow(ctx context.Context, ou org_users.OrgUser, searchGroupName string) error {
	tc := z.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: z.TaskQueueName,
		ID:        fmt.Sprintf("twitter-%s", uuid.New().String()),
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
func (z *ZeusAiPlatformServicesWorker) ExecuteAiRedditWorkflow(ctx context.Context, ou org_users.OrgUser, searchGroupName string) error {
	tc := z.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: z.TaskQueueName,
		ID:        fmt.Sprintf("reddit-%s", uuid.New().String()),
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

func (z *ZeusAiPlatformServicesWorker) ExecuteAiIngestDiscordWorkflow(ctx context.Context, ou org_users.OrgUser, cm hera_discord.ChannelMessages) error {
	tc := z.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: z.TaskQueueName,
		ID:        fmt.Sprintf("discord-%s", uuid.New().String()),
	}
	txWf := NewZeusPlatformServiceWorkflows()
	wf := txWf.AiIngestDiscordWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID, ou, cm)
	if err != nil {
		log.Err(err).Msg("ExecuteAiIngestDiscordWorkflow")
		return err
	}
	return nil
}

func (z *ZeusAiPlatformServicesWorker) ExecuteAiFetchDataToIngestDiscordWorkflow(ctx context.Context, ou org_users.OrgUser, searchGroupName string) error {
	tc := z.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: z.TaskQueueName,
		ID:        fmt.Sprintf("discord-%s", uuid.New().String()),
	}
	txWf := NewZeusPlatformServiceWorkflows()
	wf := txWf.AiFetchDataToIngestDiscordWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID, ou, searchGroupName)
	if err != nil {
		log.Err(err).Msg("AiFetchDataToIngestDiscordWorkflow")
		return err
	}
	return nil
}
