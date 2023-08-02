package quicknode_orchestrations

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_quicknode "github.com/zeus-fyi/olympus/pkg/hestia/platform/quiknode"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"go.temporal.io/sdk/workflow"
)

type HestiaQuickNodeWorkflow struct {
	temporal_base.Workflow
	HestiaQuicknodeActivities
}

const defaultTimeout = 72 * time.Hour

func NewHestiaQuickNodeWorkflow() HestiaQuickNodeWorkflow {
	deployWf := HestiaQuickNodeWorkflow{
		Workflow: temporal_base.Workflow{},
	}
	return deployWf
}

func (h *HestiaQuickNodeWorkflow) GetWorkflows() []interface{} {
	return []interface{}{h.ProvisionWorkflow, h.UpdateProvisionWorkflow, h.DeprovisionWorkflow, h.DeactivateWorkflow}
}

func (h *HestiaQuickNodeWorkflow) ProvisionWorkflow(ctx workflow.Context, ou org_users.OrgUser, pr hestia_quicknode.ProvisionRequest, user hestia_quicknode.QuickNodeUserInfo) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	pCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(pCtx, h.Provision, ou, pr, user).Get(pCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
		logger.Warn("ou", ou)
		logger.Error("failed to provision QuickNode services", "Error", err)
		return err
	}

	if !user.Verified {
		apiTokenCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(apiTokenCtx, h.InsertQuickNodeApiKey, pr).Get(apiTokenCtx, nil)
		if err != nil {
			logger.Warn("params", pr)
			logger.Warn("ou", ou)
			logger.Error("failed to provision QuickNode api for user", "Error", err)
			return err
		}
	}
	return nil
}

func (h *HestiaQuickNodeWorkflow) UpdateProvisionWorkflow(ctx workflow.Context, ou org_users.OrgUser, pr hestia_quicknode.ProvisionRequest) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	pCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(pCtx, h.UpdateProvision, pr).Get(pCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
		logger.Warn("ou", ou)
		logger.Error("failed to provision QuickNode services", "Error", err)
		return err
	}
	var excessGroups []string
	oCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(oCtx, h.CheckPlanOverages, pr).Get(oCtx, &excessGroups)
	if err != nil {
		logger.Warn("params", pr)
		logger.Warn("ou", ou)
		logger.Error("failed to adjust services", "Error", err)
		return err
	}

	if len(excessGroups) == 0 {
		return nil
	}

	for _, groupName := range excessGroups {
		dCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(pCtx, h.DeleteOrgGroupRoutingTable, ou, groupName).Get(dCtx, &excessGroups)
		if err != nil {
			logger.Warn("params", pr)
			logger.Warn("ou", ou)
			logger.Error("failed to adjust services", "Error", err)
			return err
		}
		cdCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(cdCtx, h.IrisPlatformDeleteGroupTableCacheRequest, ou, groupName).Get(cdCtx, &excessGroups)
		if err != nil {
			logger.Warn("params", pr)
			logger.Warn("ou", ou)
			logger.Error("failed to adjust cache services", "Error", err)
			return err
		}
	}
	return nil
}

func (h *HestiaQuickNodeWorkflow) DeprovisionWorkflow(ctx workflow.Context, ou org_users.OrgUser, dp hestia_quicknode.DeprovisionRequest) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	pCtx := workflow.WithActivityOptions(ctx, ao)
	currentTime := time.Now().Unix()  // get current Unix timestamp
	deprovisionAt := dp.DeprovisionAt // get provisionedAt Unix timestamp

	if currentTime < deprovisionAt {
		sleepDuration := time.Duration(deprovisionAt-currentTime) * time.Second
		err := workflow.Sleep(pCtx, sleepDuration)
		if err != nil {
			logger.Error("HestiaQuickNodeWorkflow: failed to sleep", "Error", err)
			return err
		}
	}
	err := workflow.ExecuteActivity(pCtx, h.Deprovision, ou, dp).Get(pCtx, nil)
	if err != nil {
		logger.Warn("params", dp)
		logger.Warn("ou", ou)
		logger.Error("HestiaQuickNodeWorkflow: failed to deprovision QuickNode services", "Error", err)
		return err
	}

	err = workflow.ExecuteActivity(pCtx, h.DeprovisionCache, ou).Get(pCtx, nil)
	if err != nil {
		logger.Warn("params", dp)
		logger.Warn("ou", ou)
		logger.Error("HestiaQuickNodeWorkflow: failed to DeprovisionCache", "Error", err)
		return err
	}

	err = workflow.ExecuteActivity(pCtx, h.DeactivateApiKey, ou, dp).Get(pCtx, nil)
	if err != nil {
		logger.Warn("params", dp)
		logger.Warn("ou", ou)
		logger.Error("HestiaQuickNodeWorkflow: failed to deactivate api key", "Error", err)
	}
	return nil
}

func (h *HestiaQuickNodeWorkflow) DeactivateWorkflow(ctx workflow.Context, ou org_users.OrgUser, da hestia_quicknode.DeactivateRequest) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	pCtx := workflow.WithActivityOptions(ctx, ao)
	currentTime := time.Now().Unix() // get current Unix timestamp
	deactivateAt := da.DeactivateAt  // get provisionedAt Unix timestamp

	if currentTime < deactivateAt {
		sleepDuration := time.Duration(deactivateAt-currentTime) * time.Second
		err := workflow.Sleep(pCtx, sleepDuration)
		if err != nil {
			logger.Error("HestiaQuickNodeWorkflow: failed to sleep", "Error", err)
			return err
		}
	}

	err := workflow.ExecuteActivity(pCtx, h.Deactivate, ou, da).Get(pCtx, nil)
	if err != nil {
		logger.Warn("params", da)
		logger.Warn("ou", ou)
		logger.Error("HestiaQuickNodeWorkflow: failed to deactivate QuickNode services", "Error", err)
		return err
	}
	return nil
}
