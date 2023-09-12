package quicknode_orchestrations

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_quicknode "github.com/zeus-fyi/olympus/pkg/hestia/platform/quiknode"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type HestiaQuickNodeWorkflow struct {
	temporal_base.Workflow
	HestiaQuickNodeActivities
}

const (
	defaultTimeout = 72 * time.Hour
	internalOrgID  = 7138983863666903883
)

func NewHestiaQuickNodeWorkflow() HestiaQuickNodeWorkflow {
	deployWf := HestiaQuickNodeWorkflow{
		Workflow: temporal_base.Workflow{},
	}
	return deployWf
}

func (h *HestiaQuickNodeWorkflow) GetWorkflows() []interface{} {
	return []interface{}{h.ProvisionWorkflow, h.UpdateProvisionWorkflow, h.DeprovisionWorkflow, h.DeactivateWorkflow, h.DeleteSessionCacheWorkflow}
}

func (h *HestiaQuickNodeWorkflow) ProvisionWorkflow(ctx workflow.Context, wfID string, ou org_users.OrgUser, pr hestia_quicknode.ProvisionRequest, user hestia_quicknode.QuickNodeUserInfo) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}

	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(internalOrgID, wfID, "HestiaQuickNodeWorkflow", "ProvisionWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
		logger.Warn("ou", ou)
		logger.Error("failed to provision QuickNode services", "Error", err)
		return err
	}

	pCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(pCtx, h.Provision, ou, pr, user).Get(pCtx, nil)
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
	upsertCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(upsertCtx, h.UpsertQuickNodeRoutingEndpoint, pr).Get(upsertCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
		logger.Warn("ou", ou)
		logger.Error("UpsertQuickNodeRoutingEndpoint: failed to upsert endpoint into org routing table", "Error", err)
		return err
	}
	upsertGroupCtx := workflow.WithActivityOptions(ctx, ao)
	var orgID int
	err = workflow.ExecuteActivity(upsertGroupCtx, h.UpsertQuickNodeGroupTableRoutingEndpoints, pr).Get(upsertGroupCtx, &orgID)
	if err != nil {
		logger.Warn("params", pr)
		logger.Warn("ou", ou)
		logger.Error("UpsertQuickNodeGroupTableRoutingEndpoints: failed to upsert endpoint into org routing table", "Error", err)
		return err
	}

	if orgID > 0 {
		refGroupCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(refGroupCtx, h.RefreshOrgGroupTables, orgID).Get(refGroupCtx, nil)
		if err != nil {
			logger.Warn("params", pr)
			logger.Warn("ou", ou)
			logger.Error("RefreshOrgGroupTables: failed to upsert endpoint into org routing table", "Error", err)
			return err
		}
	}

	finishedCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
		logger.Warn("ou", ou)
		logger.Error("failed to UpdateAndMarkOrchestrationInactive qn services", "Error", err)
		return err
	}
	return nil
}

func (h *HestiaQuickNodeWorkflow) DeleteSessionCacheWorkflow(ctx workflow.Context, sessionID string) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	aCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(aCtx, h.DeleteSessionAuthCache, sessionID).Get(aCtx, nil)
	if err != nil {
		return err
	}
	return nil
}

func (h *HestiaQuickNodeWorkflow) UpdateProvisionWorkflow(ctx workflow.Context, wfID string, ou org_users.OrgUser, pr hestia_quicknode.ProvisionRequest) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(internalOrgID, wfID, "HestiaQuickNodeWorkflow", "UpdateProvisionWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
		logger.Warn("ou", ou)
		logger.Error("failed to update QuickNode services", "Error", err)
		return err
	}

	pCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(pCtx, h.UpdateProvision, pr).Get(pCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
		logger.Warn("ou", ou)
		logger.Error("failed to provision QuickNode services", "Error", err)
		return err
	}
	aCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(aCtx, h.DeleteAuthCache, pr.QuickNodeID).Get(aCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
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

	for _, groupName := range excessGroups {
		dCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(pCtx, h.DeleteOrgGroupRoutingTable, ou, groupName).Get(dCtx, nil)
		if err != nil {
			logger.Warn("params", pr)
			logger.Warn("ou", ou)
			logger.Error("failed to adjust services", "Error", err)
			return err
		}
		ro := workflow.ActivityOptions{
			StartToCloseTimeout: defaultTimeout,
			RetryPolicy: &temporal.RetryPolicy{
				InitialInterval:    5 * time.Minute,
				BackoffCoefficient: 2,
				MaximumInterval:    2 * time.Minute,
			},
		}
		cdCtx := workflow.WithActivityOptions(ctx, ro)
		err = workflow.ExecuteActivity(cdCtx, h.IrisPlatformDeleteGroupTableCacheRequest, ou, groupName).Get(cdCtx, nil)
		if err != nil {
			logger.Warn("params", pr)
			logger.Warn("ou", ou)
			logger.Error("failed to adjust cache services", "Error", err)
			return err
		}
	}

	finishedCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
		logger.Warn("ou", ou)
		logger.Error("failed to UpdateAndMarkOrchestrationInactive qn services", "Error", err)
		return err
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

	deactivateKeyCtx := workflow.WithActivityOptions(ctx, ao)
	orgID := 0
	err = workflow.ExecuteActivity(deactivateKeyCtx, h.DeactivateApiKey, ou, dp).Get(deactivateKeyCtx, &orgID)
	if err != nil {
		logger.Warn("params", dp)
		logger.Warn("ou", ou)
		logger.Error("HestiaQuickNodeWorkflow: failed to deactivate api key", "Error", err)
	}
	if ou.OrgID == 0 && orgID != 0 {
		ou.OrgID = orgID
	}
	if ou.OrgID == 0 {
		log.Warn().Msg("HestiaQuickNodeWorkflow: failed to deactivate api key")
		return nil
	}
	err = workflow.ExecuteActivity(pCtx, h.DeprovisionCache, ou).Get(pCtx, nil)
	if err != nil {
		logger.Warn("params", dp)
		logger.Warn("ou", ou)
		logger.Error("HestiaQuickNodeWorkflow: failed to DeprovisionCache", "Error", err)
		return err
	}
	return nil
}

// DeactivateWorkflow removes just an endpoint
func (h *HestiaQuickNodeWorkflow) DeactivateWorkflow(ctx workflow.Context, ou org_users.OrgUser, da hestia_quicknode.DeactivateRequest) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout * 10,
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

	var httpURL string
	err := workflow.ExecuteActivity(pCtx, h.Deactivate, da).Get(pCtx, &httpURL)
	if err != nil {
		logger.Warn("params", da)
		logger.Warn("ou", ou)
		logger.Error("HestiaQuickNodeWorkflow: failed to deactivate QuickNode services", "Error", err)
		return err
	}
	if len(httpURL) == 0 {
		return nil
	}
	rmEndpointCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(rmEndpointCtx, h.IrisPlatformDeleteEndpointRequest, ou, httpURL).Get(rmEndpointCtx, nil)
	if err != nil {
		logger.Warn("params", da)
		logger.Warn("ou", ou)
		logger.Error("HestiaQuickNodeWorkflow: failed to call IrisPlatformDeleteEndpointRequest", "Error", err)
		return err
	}

	cacheRefreshCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(rmEndpointCtx, h.RefreshOrgGroupTables, ou.OrgID).Get(cacheRefreshCtx, nil)
	if err != nil {
		logger.Warn("params", da)
		logger.Warn("ou", ou)
		logger.Error("HestiaQuickNodeWorkflow: failed to call RefreshOrgGroupTables", "Error", err)
		return err
	}
	return nil
}
