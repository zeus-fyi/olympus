package platform_service_orchestrations

import (
	"errors"
	"time"

	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
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
	return []interface{}{h.IrisRoutingServiceRequestWorkflow, h.IrisDeleteOrgGroupRoutingTableWorkflow, h.IrisDeleteOrgRoutesWorkflow}
}

func (h *HestiaPlatformServiceWorkflows) IrisRoutingServiceRequestWorkflow(ctx workflow.Context, pr IrisPlatformServiceRequest) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{}
	pCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(pCtx, h.UpdateDatabaseOrgRoutingTables, pr).Get(pCtx, nil)
	if err != nil {
		log.Warn("params", pr)
		log.Error("HestiaPlatformServiceWorkflows: failed to UpdateDatabaseOrgRoutingTables", "Error", err)
		return err
	}
	if pr.OrgGroupName != "" {
		err = workflow.ExecuteActivity(pCtx, h.CreateOrgGroupRoutingTable, pr).Get(pCtx, nil)
		if err != nil {
			log.Warn("params", pr)
			log.Error("HestiaPlatformServiceWorkflows: failed to CreateOrgGroupRoutingTable", "Error", err)
			return err
		}
	}
	err = workflow.ExecuteActivity(pCtx, h.IrisPlatformSetupCacheUpdateRequest, pr).Get(pCtx, nil)
	if err != nil {
		log.Warn("params", pr)
		log.Error("HestiaPlatformServiceWorkflows: failed to complete IrisPlatformSetupCacheUpdateRequest", "Error", err)
		return err
	}
	return nil
}

func (h *HestiaPlatformServiceWorkflows) IrisDeleteOrgGroupRoutingTableWorkflow(ctx workflow.Context, pr IrisPlatformServiceRequest) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{}
	if pr.OrgGroupName == "" {
		return errors.New("HestiaPlatformServiceWorkflows: IrisDeleteOrgGroupRoutingTableWorkflow: org group name is empty")
	}
	pCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(pCtx, h.DeleteOrgGroupRoutingTable, pr).Get(pCtx, nil)
	if err != nil {
		log.Warn("params", pr)
		log.Error("HestiaPlatformServiceWorkflows: IrisDeleteOrgGroupRoutingTableWorkflow: failed to DeleteOrgGroupRoutingTable", "Error", err)
		return err
	}
	err = workflow.ExecuteActivity(pCtx, h.IrisPlatformDeleteGroupTableCacheRequest, pr).Get(pCtx, nil)
	if err != nil {
		log.Warn("params", pr)
		log.Error("HestiaPlatformServiceWorkflows: failed to complete IrisPlatformDeleteGroupTableCacheRequest", "Error", err)
		return err
	}
	return nil
}

func (h *HestiaPlatformServiceWorkflows) IrisDeleteOrgRoutesWorkflow(ctx workflow.Context, pr IrisPlatformServiceRequest) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{}
	if len(pr.Routes) == 0 {
		return errors.New("HestiaPlatformServiceWorkflows: IrisDeleteOrgRoutesWorkflow: no routes provided for deletion")
	}
	pCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(pCtx, h.DeleteOrgRoutes, pr).Get(pCtx, nil)
	if err != nil {
		log.Warn("params", pr)
		log.Error("HestiaPlatformServiceWorkflows: failed to DeleteOrgRoutes", "Error", err)
		return err
	}
	err = workflow.ExecuteActivity(pCtx, h.IrisPlatformSetupCacheUpdateRequest, pr).Get(pCtx, nil)
	if err != nil {
		log.Warn("params", pr)
		log.Error("HestiaPlatformServiceWorkflows: failed to complete IrisPlatformSetupCacheUpdateRequest", "Error", err)
		return err
	}
	return nil
}

func (h *HestiaPlatformServiceWorkflows) IrisRemoveAllOrgRoutesFromCacheWorkflow(ctx workflow.Context, pr IrisPlatformServiceRequest) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{}
	pCtx := workflow.WithActivityOptions(ctx, ao)
	// this deletes all the routing tables from cache for this org
	err := workflow.ExecuteActivity(pCtx, h.IrisPlatformDeleteOrgGroupTablesCacheRequest, pr).Get(pCtx, nil)
	if err != nil {
		log.Warn("params", pr)
		log.Error("HestiaPlatformServiceWorkflows: failed to complete IrisPlatformDeleteOrgGroupTablesCacheRequest", "Error", err)
		return err
	}
	return nil
}
