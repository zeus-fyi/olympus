package quicknode_orchestrations

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	hestia_quicknode "github.com/zeus-fyi/olympus/pkg/hestia/platform/quiknode"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"go.temporal.io/sdk/workflow"
)

type HestiaQuicknodeWorkflow struct {
	temporal_base.Workflow
	HestiaQuicknodeActivities
}

const defaultTimeout = 72 * time.Hour

func NewHestiaQuicknodeWorkflow() HestiaQuicknodeWorkflow {
	deployWf := HestiaQuicknodeWorkflow{
		Workflow: temporal_base.Workflow{},
	}
	return deployWf
}

func (h *HestiaQuicknodeWorkflow) GetWorkflows() []interface{} {
	return []interface{}{h.ProvisionWorkflow, h.UpdateProvisionWorkflow, h.DeprovisionWorkflow, h.DeactivateWorkflow,
		h.DeleteOrgGroupRoutingTable}
}

func (h *HestiaQuicknodeWorkflow) DeleteOrgGroupRoutingTable(ctx context.Context, ou org_users.OrgUser, groupName string) error {
	err := iris_models.DeleteOrgGroupAndRoutes(context.Background(), ou.OrgID, groupName)
	if err != nil {
		log.Err(err).Msg("DeleteOrgGroupRoutingTable: DeleteOrgGroupRoutingTable")
		return err
	}
	return nil
}

func (h *HestiaQuicknodeWorkflow) ProvisionWorkflow(ctx workflow.Context, pr hestia_quicknode.ProvisionRequest, ou org_users.OrgUser, user hestia_quicknode.QuickNodeUserInfo) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	pCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(pCtx, h.Provision, pr, ou, user).Get(pCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
		logger.Warn("ou", ou)
		logger.Error("failed to provision QuickNode services", "Error", err)
		return err
	}
	return nil
}

func (h *HestiaQuicknodeWorkflow) UpdateProvisionWorkflow(ctx workflow.Context, pr hestia_quicknode.ProvisionRequest, ou org_users.OrgUser) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	pCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(pCtx, h.UpdateProvision, pr, ou).Get(pCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
		logger.Warn("ou", ou)
		logger.Error("failed to provision QuickNode services", "Error", err)
		return err
	}
	var excessGroups []string
	err = workflow.ExecuteActivity(pCtx, h.CheckPlanOverages, pr, ou).Get(pCtx, &excessGroups)
	if err != nil {
		logger.Warn("params", pr)
		logger.Warn("ou", ou)
		logger.Error("failed to adjust services", "Error", err)
		return err
	}

	for _, groupName := range excessGroups {
		err = workflow.ExecuteActivity(pCtx, h.DeleteOrgGroupRoutingTable, ou, groupName).Get(pCtx, &excessGroups)
		if err != nil {
			logger.Warn("params", pr)
			logger.Warn("ou", ou)
			logger.Error("failed to adjust services", "Error", err)
			return err
		}
		err = workflow.ExecuteActivity(pCtx, h.IrisPlatformDeleteGroupTableCacheRequest, ou, groupName).Get(pCtx, &excessGroups)
		if err != nil {
			logger.Warn("params", pr)
			logger.Warn("ou", ou)
			logger.Error("failed to adjust cache services", "Error", err)
			return err
		}
	}
	return nil
}

func (h *HestiaQuicknodeWorkflow) DeprovisionWorkflow(ctx workflow.Context, dp hestia_quicknode.DeprovisionRequest, ou org_users.OrgUser) error {
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
			logger.Error("HestiaQuicknodeWorkflow: failed to sleep", "Error", err)
			return err
		}
	}
	err := workflow.ExecuteActivity(pCtx, h.Deprovision, dp, ou).Get(pCtx, nil)
	if err != nil {
		logger.Warn("params", dp)
		logger.Warn("ou", ou)
		logger.Error("HestiaQuicknodeWorkflow: failed to deprovision QuickNode services", "Error", err)
		return err
	}

	err = workflow.ExecuteActivity(pCtx, h.DeprovisionCache, ou).Get(pCtx, nil)
	if err != nil {
		logger.Warn("params", dp)
		logger.Warn("ou", ou)
		logger.Error("HestiaQuicknodeWorkflow: failed to DeprovisionCache", "Error", err)
		return err
	}
	return nil
}

func (h *HestiaQuicknodeWorkflow) DeactivateWorkflow(ctx workflow.Context, da hestia_quicknode.DeactivateRequest, ou org_users.OrgUser) error {
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
			logger.Error("HestiaQuicknodeWorkflow: failed to sleep", "Error", err)
			return err
		}
	}

	err := workflow.ExecuteActivity(pCtx, h.Deactivate, da, ou).Get(pCtx, nil)
	if err != nil {
		logger.Warn("params", da)
		logger.Warn("ou", ou)
		logger.Error("HestiaQuicknodeWorkflow: failed to deactivate QuickNode services", "Error", err)
		return err
	}
	return nil
}
