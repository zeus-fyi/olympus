package deploy_workflow_destroy_setup

import (
	"context"
	"time"

	hestia_stripe "github.com/zeus-fyi/olympus/pkg/hestia/stripe"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	deploy_topology_activities_create_setup "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/cluster_setup"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"go.temporal.io/sdk/workflow"
)

type DestroyClusterSetupWorkflow struct {
	temporal_base.Workflow
	deploy_topology_activities_create_setup.CreateSetupTopologyActivities
}

const defaultTimeout = 120 * time.Minute

func NewDeployDestroyClusterSetupWorkflow() DestroyClusterSetupWorkflow {
	deployWf := DestroyClusterSetupWorkflow{
		Workflow:                      temporal_base.Workflow{},
		CreateSetupTopologyActivities: deploy_topology_activities_create_setup.CreateSetupTopologyActivities{},
	}
	return deployWf
}

func (c *DestroyClusterSetupWorkflow) GetWorkflow() interface{} {
	return c.DestroyClusterSetupWorkflow
}

func (c *DestroyClusterSetupWorkflow) GetWorkflows() []interface{} {
	return []interface{}{c.DestroyClusterSetupWorkflow}
}

func (c *DestroyClusterSetupWorkflow) DestroyClusterSetupWorkflow(ctx workflow.Context, params base_deploy_params.DestroyClusterSetupRequest) error {
	log := workflow.GetLogger(ctx)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	if params.FreeTrial {
		err := workflow.Sleep(ctx, 60*time.Minute)
		if err != nil {
			log.Error("Failed to sleep for 1 hour", "Error", err)
			return err
		}
		hestiaCtx := context.Background()
		isBillingSetup, herr := hestia_stripe.DoesUserHaveBillingMethod(hestiaCtx, params.Ou.UserID)
		if herr != nil {
			log.Error("Failed to check if user has billing method", "Error", herr)
			return herr
		}
		if !isBillingSetup {
			destroyClusterCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(destroyClusterCtx, c.CreateSetupTopologyActivities.DestroyCluster, params.CloudCtxNs).Get(destroyClusterCtx, nil)
			if err != nil {
				log.Error("Failed to add deploy cluster", "Error", err)
				return err
			}
			removeAuthCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(removeAuthCtx, c.CreateSetupTopologyActivities.RemoveAuthCtxNsOrg, params).Get(removeAuthCtx, nil)
			if err != nil {
				log.Error("Failed to add deploy cluster", "Error", err)
				return err
			}
			destroyNodePoolOrgResourcesCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(destroyNodePoolOrgResourcesCtx, c.CreateSetupTopologyActivities.RemoveNodePoolRequest, params.DigitalOceanNodePoolRequestStatus).Get(destroyNodePoolOrgResourcesCtx, nil)
			if err != nil {
				log.Error("Failed to add remove node resources for account", "Error", err)
				return err
			}
			removeFreeTrialResourcesCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(removeFreeTrialResourcesCtx, c.CreateSetupTopologyActivities.RemoveFreeTrialOrgResources, params).Get(removeFreeTrialResourcesCtx, nil)
			if err != nil {
				log.Error("Failed to add remove org free trial resources for account", "Error", err)
				return err
			}
		} else {
			updateResourcesToPaidCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(updateResourcesToPaidCtx, c.CreateSetupTopologyActivities.UpdateFreeTrialOrgResourcesToPaid, params).Get(updateResourcesToPaidCtx, nil)
			if err != nil {
				log.Error("Failed to update org free trial resources to paid for account", "Error", err)
				return err
			}
		}
	}
	return nil
}
