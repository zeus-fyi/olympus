package platform_service_orchestrations

import (
	"errors"
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type HestiaPlatformServiceWorkflows struct {
	temporal_base.Workflow
	HestiaPlatformActivities
}

const defaultTimeout = 72 * time.Hour

func NewHestiaPlatformServiceWorkflows() HestiaPlatformServiceWorkflows {
	deployWf := HestiaPlatformServiceWorkflows{
		Workflow: temporal_base.Workflow{},
	}
	return deployWf
}

func (h *HestiaPlatformServiceWorkflows) GetWorkflows() []interface{} {
	return []interface{}{h.IrisRoutingServiceRequestWorkflow, h.IrisDeleteOrgGroupRoutingTableWorkflow, h.IrisDeleteOrgRoutesWorkflow,
		h.IrisRemoveAllOrgRoutesFromCacheWorkflow, h.IrisDeleteRoutesFromOrgGroupRoutingTableWorkflow}
}

const (
	internalOrgID = 7138983863666903883
)

func (h *HestiaPlatformServiceWorkflows) IrisRoutingServiceRequestWorkflow(ctx workflow.Context, wfID string, pr IrisPlatformServiceRequest) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
	}
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(internalOrgID, wfID, "HestiaPlatformServiceWorkflows", "IrisRoutingServiceRequestWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
		logger.Warn("ou", pr.Ou)
		logger.Error("failed to update QuickNode services", "Error", err)
		return err
	}

	pCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(pCtx, h.UpdateDatabaseOrgRoutingTables, pr).Get(pCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
		logger.Error("HestiaPlatformServiceWorkflows: failed to UpdateDatabaseOrgRoutingTables", "Error", err)
		return err
	}
	if pr.OrgGroupName != "" {
		err = workflow.ExecuteActivity(pCtx, h.CreateOrgGroupRoutingTable, pr).Get(pCtx, nil)
		if err != nil {
			logger.Warn("params", pr)
			logger.Error("HestiaPlatformServiceWorkflows: failed to CreateOrgGroupRoutingTable", "Error", err)
			return err
		}
		so := workflow.ActivityOptions{
			StartToCloseTimeout: time.Minute * 15,
			RetryPolicy: &temporal.RetryPolicy{
				InitialInterval:    time.Minute * 1,
				BackoffCoefficient: 1.5,
				MaximumInterval:    time.Minute * 5,
			},
		}
		sCtx := workflow.WithActivityOptions(ctx, so)
		err = workflow.ExecuteActivity(sCtx, h.IrisPlatformSetupCacheUpdateRequest, pr).Get(sCtx, nil)
		if err != nil {
			logger.Warn("params", pr)
			logger.Error("HestiaPlatformServiceWorkflows: failed to complete IrisPlatformSetupCacheUpdateRequest", "Error", err)
			return err
		}
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

func (h *HestiaPlatformServiceWorkflows) IrisDeleteOrgGroupRoutingTableWorkflow(ctx workflow.Context, wfID string, pr IrisPlatformServiceRequest) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 15,
	}
	if pr.OrgGroupName == "" {
		return errors.New("HestiaPlatformServiceWorkflows: IrisDeleteOrgGroupRoutingTableWorkflow: org group name is empty")
	}

	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(internalOrgID, wfID, "HestiaPlatformServiceWorkflows", "IrisDeleteOrgGroupRoutingTableWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
		logger.Warn("ou", pr.Ou)
		logger.Error("failed to update cache for qn services", "Error", err)
		return err
	}

	orgGroupName := pr.OrgGroupName
	pCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(pCtx, h.DeleteOrgRoutingGroup, pr).Get(pCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
		logger.Error("HestiaPlatformServiceWorkflows: DeleteOrgRoutesFromGroup: failed to DeleteOrgGroupRoutingTable", "Error", err)
		return err
	}
	pr.OrgGroupName = orgGroupName
	do := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 15,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Minute * 1,
			BackoffCoefficient: 1.5,
			MaximumInterval:    time.Minute * 5,
		},
	}
	dCtx := workflow.WithActivityOptions(ctx, do)
	err = workflow.ExecuteActivity(dCtx, h.IrisPlatformDeleteGroupTableCacheRequest, pr).Get(dCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
		logger.Error("HestiaPlatformServiceWorkflows: failed to complete IrisPlatformDeleteGroupTableCacheRequest", "Error", err)
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

func (h *HestiaPlatformServiceWorkflows) IrisDeleteRoutesFromOrgGroupRoutingTableWorkflow(ctx workflow.Context, wfID string, pr IrisPlatformServiceRequest) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 30,
	}
	if pr.OrgGroupName == "" {
		return errors.New("HestiaPlatformServiceWorkflows: IrisDeleteOrgGroupRoutingTableWorkflow: org group name is empty")
	}
	if len(pr.Routes) == 0 {
		return errors.New("HestiaPlatformServiceWorkflows: IrisDeleteOrgGroupRoutingTableWorkflow: no routes provided for deletion")
	}

	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(internalOrgID, wfID, "HestiaPlatformServiceWorkflows", "IrisDeleteRoutesFromOrgGroupRoutingTableWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
		logger.Warn("ou", pr.Ou)
		logger.Error("failed to update cache for qn services", "Error", err)
		return err
	}

	orgGroupName := pr.OrgGroupName
	pCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(pCtx, h.DeleteOrgRoutesFromGroup, pr).Get(pCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
		logger.Error("HestiaPlatformServiceWorkflows: DeleteOrgRoutesFromGroup: failed to DeleteOrgGroupRoutingTable", "Error", err)
		return err
	}
	pr.OrgGroupName = orgGroupName
	do := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 30,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Minute * 1,
			BackoffCoefficient: 1.5,
			MaximumInterval:    time.Minute * 5,
		},
	}
	dCtx := workflow.WithActivityOptions(ctx, do)
	err = workflow.ExecuteActivity(dCtx, h.IrisPlatformRefreshOrgGroupTableCacheRequest, pr).Get(dCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
		logger.Error("HestiaPlatformServiceWorkflows: failed to complete IrisPlatformRefreshOrgGroupTableCacheRequest", "Error", err)
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

func (h *HestiaPlatformServiceWorkflows) IrisDeleteOrgRoutesWorkflow(ctx workflow.Context, wfID string, pr IrisPlatformServiceRequest) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 30,
	}
	if len(pr.Routes) == 0 {
		return errors.New("HestiaPlatformServiceWorkflows: IrisDeleteOrgRoutesWorkflow: no routes provided for deletion")
	}
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(internalOrgID, wfID, "HestiaPlatformServiceWorkflows", "IrisDeleteOrgRoutesWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
		logger.Warn("ou", pr.Ou)
		logger.Error("failed to update QuickNode services", "Error", err)
		return err
	}

	pCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(pCtx, h.DeleteOrgRoutes, pr).Get(pCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
		logger.Error("HestiaPlatformServiceWorkflows: failed to DeleteOrgRoutes", "Error", err)
		return err
	}
	do := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 15,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Minute * 1,
			BackoffCoefficient: 1.5,
			MaximumInterval:    time.Minute * 5,
		},
	}
	dCtx := workflow.WithActivityOptions(ctx, do)
	err = workflow.ExecuteActivity(dCtx, h.IrisPlatformSetupCacheUpdateRequest, pr).Get(dCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
		logger.Error("HestiaPlatformServiceWorkflows: failed to complete IrisPlatformSetupCacheUpdateRequest", "Error", err)
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

func (h *HestiaPlatformServiceWorkflows) IrisRemoveAllOrgRoutesFromCacheWorkflow(ctx workflow.Context, wfID string, pr IrisPlatformServiceRequest) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 15,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Minute * 1,
			BackoffCoefficient: 1.5,
			MaximumInterval:    time.Minute * 5,
		},
	}
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(internalOrgID, wfID, "HestiaPlatformServiceWorkflows", "IrisRemoveAllOrgRoutesFromCacheWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
		logger.Warn("ou", pr.Ou)
		logger.Error("failed to update QuickNode services", "Error", err)
		return err
	}

	pCtx := workflow.WithActivityOptions(ctx, ao)
	// this deletes all the routing tables from cache for this org
	err = workflow.ExecuteActivity(pCtx, h.IrisPlatformDeleteOrgGroupTablesCacheRequest, pr).Get(pCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
		logger.Error("HestiaPlatformServiceWorkflows: failed to complete IrisPlatformDeleteOrgGroupTablesCacheRequest", "Error", err)
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
