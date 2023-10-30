package iris_serverless

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"go.temporal.io/sdk/workflow"
)

type IrisPlatformServiceWorkflows struct {
	temporal_base.Workflow
	IrisPlatformActivities
}

const defaultTimeout = 72 * time.Hour

func NewIrisPlatformServiceWorkflows() IrisPlatformServiceWorkflows {
	deployWf := IrisPlatformServiceWorkflows{
		Workflow: temporal_base.Workflow{},
	}
	return deployWf
}

func (i *IrisPlatformServiceWorkflows) GetWorkflows() []interface{} {
	return []interface{}{i.IrisServerlessResyncWorkflow}
}

func (i *IrisPlatformServiceWorkflows) IrisServerlessResyncWorkflow(ctx workflow.Context, wfID string) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
	}
	oj := artemis_orchestrations.NewInternalActiveTemporalOrchestrationJobTemplate(wfID, "IrisPlatformServiceWorkflows", "IrisServerlessResyncWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Error("failed to update QuickNode services", "Error", err)
		return err
	}
	pCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(pCtx, i.ResyncServerlessRoutes, nil).Get(pCtx, nil)
	if err != nil {
		logger.Error("IrisPlatformServiceWorkflows: failed to ResyncServerlessRoutes", "Error", err)
		return err
	}
	finishedCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
	if err != nil {
		logger.Error("failed to update cache for qn services", "Error", err)
		return err
	}
	return nil
}