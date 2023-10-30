package platform_service_orchestrations

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"go.temporal.io/sdk/workflow"
)

func (h *HestiaPlatformServiceWorkflows) IrisServerlessResyncWorkflow(ctx workflow.Context, wfID string, pr IrisPlatformServiceRequest) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
	}
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(internalOrgID, wfID, "HestiaPlatformServiceWorkflows", "IrisServerlessResyncWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
		logger.Warn("ou", pr.Ou)
		logger.Error("failed to update QuickNode services", "Error", err)
		return err
	}

	pCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(pCtx, h.ResyncServerlessRoutes, pr).Get(pCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
		logger.Error("HestiaPlatformServiceWorkflows: failed to ResyncServerlessRoutes", "Error", err)
		return err
	}

	finishedCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
		logger.Warn("ou", pr.Ou)
		logger.Error("failed to update cache for qn services", "Error", err)
		return err
	}
	return nil
}
