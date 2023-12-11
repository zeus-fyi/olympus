package ai_platform_service_orchestrations

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"go.temporal.io/sdk/workflow"
)

func (z *ZeusAiPlatformServiceWorkflows) CancelWorkflowRuns(ctx workflow.Context, wfID string, ou org_users.OrgUser, wfIDs []string) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 15, // Setting a valid non-zero timeout
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	for _, runID := range wfIDs {
		err := workflow.ExecuteActivity(ctx, z.CancelRun, runID).Get(ctx, nil)
		if err != nil {
			logger.Error("CancelRunsWorkflow: ExecuteActivity failed.", "Error", err)
			return err
		}
		finishedCtx := workflow.WithActivityOptions(ctx, ao)
		oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(ou.OrgID, runID, "", "")
		err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
		if err != nil {
			logger.Error("failed to update orchestration status", "Error", err)
			return err
		}
	}
	return nil
}
