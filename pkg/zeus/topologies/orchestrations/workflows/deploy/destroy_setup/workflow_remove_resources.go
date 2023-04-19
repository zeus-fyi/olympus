package deploy_workflow_destroy_setup

import (
	do_types "github.com/zeus-fyi/olympus/pkg/hestia/digitalocean/types"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	deploy_topology_activities_create_setup "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/cluster_setup"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"go.temporal.io/sdk/workflow"
)

type DestroyResourcesWorkflow struct {
	temporal_base.Workflow
	deploy_topology_activities_create_setup.CreateSetupTopologyActivities
}

func NewDestroyResourcesWorkflow() DestroyResourcesWorkflow {
	deployWf := DestroyResourcesWorkflow{
		Workflow:                      temporal_base.Workflow{},
		CreateSetupTopologyActivities: deploy_topology_activities_create_setup.CreateSetupTopologyActivities{},
	}
	return deployWf
}

func (c *DestroyResourcesWorkflow) GetWorkflow() interface{} {
	return c.DestroyClusterResourcesWorkflow
}

func (c *DestroyResourcesWorkflow) GetWorkflows() []interface{} {
	return []interface{}{c.DestroyClusterResourcesWorkflow}
}

func (c *DestroyResourcesWorkflow) DestroyClusterResourcesWorkflow(ctx workflow.Context, params base_deploy_params.DestroyResourcesRequest) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	selectNodesCtx := workflow.WithActivityOptions(ctx, ao)
	var nodes []do_types.DigitalOceanNodePoolRequestStatus
	err := workflow.ExecuteActivity(selectNodesCtx, c.CreateSetupTopologyActivities.SelectNodeResources, params).Get(selectNodesCtx, &nodes)
	if err != nil {
		log.Error("Failed to select org resource nodes", "Error", err)
		return err
	}
	if len(nodes) == 0 {
		log.Info("No node resources found to destroy or they were free trial nodes that will be deleted automatically")
		return nil
	}
	for _, node := range nodes {
		destroyNodePoolOrgResourcesCtx := workflow.WithActivityOptions(ctx, ao)
		log.Info("Destroying node pool org resources", "NodePoolRequestStatus", node)
		err = workflow.ExecuteActivity(destroyNodePoolOrgResourcesCtx, c.CreateSetupTopologyActivities.RemoveNodePoolRequest, node).Get(destroyNodePoolOrgResourcesCtx, nil)
		if err != nil {
			log.Error("Failed to remove node resources for account", "Error", err)
			return err
		}
	}
	// TODO billing somewhere usage update
	endServiceNodesCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(endServiceNodesCtx, c.CreateSetupTopologyActivities.EndResourceService, params).Get(endServiceNodesCtx, &nodes)
	if err != nil {
		log.Error("Failed to update org_resources to end service", "Error", err)
		return err
	}

	return nil
}
