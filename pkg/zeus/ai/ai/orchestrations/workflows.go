package ai_platform_service_orchestrations

import (
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hermes_email_notifications "github.com/zeus-fyi/olympus/pkg/hermes/email"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type ZeusAiPlatformServiceWorkflows struct {
	temporal_base.Workflow
	ZeusAiPlatformActivities
}

const defaultTimeout = 72 * time.Hour

func NewZeusPlatformServiceWorkflows() ZeusAiPlatformServiceWorkflows {
	deployWf := ZeusAiPlatformServiceWorkflows{
		Workflow: temporal_base.Workflow{},
	}
	return deployWf
}

func (h *ZeusAiPlatformServiceWorkflows) GetWorkflows() []interface{} {
	return []interface{}{h.AiEmailWorkflow}
}

const (
	internalOrgID = 7138983863666903883
)

func (h *ZeusAiPlatformServiceWorkflows) AiEmailWorkflow(ctx workflow.Context, wfID string, ou org_users.OrgUser, msgs []hermes_email_notifications.EmailContents) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
		},
	}
	// todo allow user orgs ids
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(internalOrgID, wfID, "ZeusAiPlatformServiceWorkflows", "AiEmailWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Error("failed to update ai orch services", "Error", err)
		return err
	}

	for _, msg := range msgs {
		var emailID int
		insertEmailCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(insertEmailCtx, h.InsertEmailIfNew, msg).Get(insertEmailCtx, &emailID)
		if err != nil {
			logger.Error("failed to execute InsertEmailIfNew", "Error", err)
			// You can decide if you want to return the error or continue monitoring.
			return err
		}
		if emailID <= 0 {
			continue
		}
		var resp openai.ChatCompletionResponse
		runAiTaskCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(runAiTaskCtx, h.AiTask, ou, msg).Get(runAiTaskCtx, &resp)
		if err != nil {
			logger.Error("failed to execute AiTask", "Error", err)
			// You can decide if you want to return the error or continue monitoring.
			return err
		}

		sendEmailTaskCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(sendEmailTaskCtx, h.SendTaskResponseEmail, msg.From, resp).Get(sendEmailTaskCtx, &resp)
		if err != nil {
			logger.Error("failed to execute SaveAiTaskResponse", "Error", err)
			// You can decide if you want to return the error or continue monitoring.
			return err
		}

		if ou.OrgID > 0 && ou.UserID > 0 {
			saveAiTaskCompletionCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(saveAiTaskCompletionCtx, h.SaveAiTaskResponse, ou, resp).Get(saveAiTaskCompletionCtx, &resp)
			if err != nil {
				logger.Error("failed to execute SaveAiTaskResponse", "Error", err)
				// You can decide if you want to return the error or continue monitoring.
				return err
			}
		}
	}

	finishedCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
	if err != nil {
		logger.Error("failed to update cache for qn services", "Error", err)
		return err
	}
	return nil
}
