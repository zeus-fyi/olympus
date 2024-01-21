package ai_platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	hera_openai_dbmodels "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"go.temporal.io/sdk/client"
)

func (z *ZeusAiPlatformServicesWorker) ExecuteSocialMediaExtractionWorkflow(ctx context.Context, ou org_users.OrgUser,
	sg *hera_openai_dbmodels.SearchResultGroup) (*ChatCompletionQueryResponse, error) {
	tc := z.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: z.TaskQueueName,
		ID:        fmt.Sprintf("sm-extraction-%s", uuid.New().String()),
	}
	txWf := NewZeusPlatformServiceWorkflows()
	wf := txWf.SocialMediaExtractionWorkflow
	var cr *ChatCompletionQueryResponse
	workflowRun, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID, ou, sg)
	if err != nil {
		log.Err(err).Msg("ExecuteSmExtractionWorkflow")
		return nil, err
	}
	err = workflowRun.Get(ctx, &cr)
	if err != nil {
		log.Err(err).Msg("ExecuteSmExtractionWorkflow: Get ChatCompletionQueryResponse")
		return nil, err
	}
	return cr, nil
}

func (z *ZeusAiPlatformServicesWorker) ExecuteSocialMediaEngagementWorkflow(ctx context.Context, ou org_users.OrgUser,
	sg *hera_openai_dbmodels.SearchResultGroup) (*ChatCompletionQueryResponse, error) {
	tc := z.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: z.TaskQueueName,
		ID:        fmt.Sprintf("sm-engagement-%s", uuid.New().String()),
	}
	txWf := NewZeusPlatformServiceWorkflows()
	wf := txWf.SocialMediaEngagementWorkflow
	var cr *ChatCompletionQueryResponse
	workflowRun, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID, ou, sg)
	if err != nil {
		log.Err(err).Msg("ExecuteSocialMediaEngagementWorkflow")
		return nil, err
	}
	err = workflowRun.Get(ctx, &cr)
	if err != nil {
		log.Err(err).Msg("ExecuteSocialMediaEngagementWorkflow: Get ChatCompletionQueryResponse")
		return nil, err
	}
	return cr, nil
}
