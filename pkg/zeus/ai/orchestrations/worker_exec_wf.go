package ai_platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hermes_email_notifications "github.com/zeus-fyi/olympus/pkg/hermes/email"
	"go.temporal.io/sdk/client"
)

func (z *ZeusAiPlatformServicesWorker) ExecuteAiTaskWorkflow(ctx context.Context, ou org_users.OrgUser, msgs []hermes_email_notifications.EmailContents) error {
	tc := z.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: z.TaskQueueName,
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

func (z *ZeusAiPlatformServicesWorker) ExecuteCancelWorkflowRuns(ctx context.Context, ou org_users.OrgUser, wfIDs []string) error {
	tc := z.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: z.TaskQueueName,
		ID:        fmt.Sprintf("cancel-runs-%s", uuid.New().String()),
	}
	txWf := NewZeusPlatformServiceWorkflows()
	wf := txWf.CancelWorkflowRuns
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID, ou, wfIDs)
	if err != nil {
		log.Err(err).Msg("ExecuteCancelWorkflowRuns")
		return err
	}
	return nil
}

func (z *ZeusAiPlatformServicesWorker) ExecuteCancelWorkflow(ctx context.Context, wfID string) error {
	tc := z.ConnectTemporalClient()
	defer tc.Close()
	err := tc.TerminateWorkflow(ctx, wfID, "", "user requested")
	if err != nil {
		log.Err(err).Msg("ExecuteCancelWorkflow")
		return err
	}
	return nil
}
