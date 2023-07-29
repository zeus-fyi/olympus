package quicknode_orchestrations

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
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
	return []interface{}{h.ProvisionWorkflow, h.UpdateProvisionWorkflow, h.DeprovisionWorkflow, h.DeactivateWorkflow}
}

func (h *HestiaQuicknodeWorkflow) ProvisionWorkflow(ctx workflow.Context, pr hestia_quicknode.ProvisionRequest, ou org_users.OrgUser) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	pCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(pCtx, h.Provision, pr, ou).Get(pCtx, nil)
	if err != nil {
		log.Warn("params", pr)
		log.Warn("ou", ou)
		log.Error("failed to provision QuickNode services", "Error", err)
		return err
	}
	return nil
}

func (h *HestiaQuicknodeWorkflow) UpdateProvisionWorkflow(ctx workflow.Context, pr hestia_quicknode.ProvisionRequest, ou org_users.OrgUser) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	pCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(pCtx, h.UpdateProvision, pr, ou).Get(pCtx, nil)
	if err != nil {
		log.Warn("params", pr)
		log.Warn("ou", ou)
		log.Error("failed to provision QuickNode services", "Error", err)
		return err
	}
	return nil
}

func (h *HestiaQuicknodeWorkflow) DeprovisionWorkflow(ctx workflow.Context, dp hestia_quicknode.DeprovisionRequest, ou org_users.OrgUser) error {
	log := workflow.GetLogger(ctx)
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
			return err
		}
	}
	err := workflow.ExecuteActivity(pCtx, h.Deprovision, dp, ou).Get(pCtx, nil)
	if err != nil {
		log.Warn("params", dp)
		log.Warn("ou", ou)
		log.Error("failed to deprovision QuickNode services", "Error", err)
		return err
	}

	err = workflow.ExecuteActivity(pCtx, h.DeprovisionCache, ou).Get(pCtx, nil)
	if err != nil {
		log.Warn("params", dp)
		log.Warn("ou", ou)
		log.Error("failed to DeprovisionCache", "Error", err)
		return err
	}
	return nil
}

func (h *HestiaQuicknodeWorkflow) DeactivateWorkflow(ctx workflow.Context, da hestia_quicknode.DeactivateRequest, ou org_users.OrgUser) error {
	log := workflow.GetLogger(ctx)
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
			return err
		}
	}

	err := workflow.ExecuteActivity(pCtx, h.Deactivate, da, ou).Get(pCtx, nil)
	if err != nil {
		log.Warn("params", da)
		log.Warn("ou", ou)
		log.Error("failed to deactivate QuickNode services", "Error", err)
		return err
	}
	return nil
}
