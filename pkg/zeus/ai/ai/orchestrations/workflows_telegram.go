package ai_platform_service_orchestrations

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_openai_dbmodels "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (h *ZeusAiPlatformServiceWorkflows) AiIngestTelegramWorkflow(ctx workflow.Context, wfID string, ou org_users.OrgUser, msgs []hera_openai_dbmodels.TelegramMessage) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
		},
	}
	// todo allow user orgs ids
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(internalOrgID, wfID, "ZeusAiPlatformServiceWorkflows", "AiIngestTelegramWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Error("failed to update ai orch services", "Error", err)
		return err
	}

	for _, msg := range msgs {
		var msgID int
		insertMsgCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(insertMsgCtx, h.InsertTelegramMessageIfNew, ou, msg).Get(insertMsgCtx, &msgID)
		if err != nil {
			logger.Error("failed to execute InsertEmailIfNew", "Error", err)
			// You can decide if you want to return the error or continue monitoring.
			return err
		}
		if msgID <= 0 {
			continue
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
