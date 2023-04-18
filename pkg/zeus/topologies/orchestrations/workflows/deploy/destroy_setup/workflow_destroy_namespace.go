package deploy_workflow_destroy_setup

import (
	hestia_compute_resources "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/resources"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	deploy_topology_activities_create_setup "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/cluster_setup"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

type DestroyNamespaceSetupWorkflow struct {
	temporal_base.Workflow
	deploy_topology_activities_create_setup.CreateSetupTopologyActivities
}

func NewDestroyNamespaceSetupWorkflow() DestroyNamespaceSetupWorkflow {
	deployWf := DestroyNamespaceSetupWorkflow{
		Workflow:                      temporal_base.Workflow{},
		CreateSetupTopologyActivities: deploy_topology_activities_create_setup.CreateSetupTopologyActivities{},
	}
	return deployWf
}

func (c *DestroyNamespaceSetupWorkflow) GetWorkflow() interface{} {
	return c.DestroyNamespaceSetupWorkflow
}

func (c *DestroyNamespaceSetupWorkflow) GetWorkflows() []interface{} {
	return []interface{}{c.DestroyNamespaceSetupWorkflow}
}

func (c *DestroyNamespaceSetupWorkflow) DestroyNamespaceSetupWorkflow(ctx workflow.Context, params base_deploy_params.TopologyWorkflowRequest) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	removeAuthCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(removeAuthCtx, c.CreateSetupTopologyActivities.RemoveAuthCtxNsOrg, params).Get(removeAuthCtx, nil)
	if err != nil {
		log.Error("Failed to remove auth ctx ns", "Error", err)
		return err
	}
	removeSubdomainCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(removeSubdomainCtx, c.CreateSetupTopologyActivities.RemoveDomainRecord, params.Kns.CloudCtxNs.Namespace).Get(removeSubdomainCtx, nil)
	if err != nil {
		log.Error("Failed to remove domain record", "Error", err)
		return err
	}
	getDisksAtCloudCtxNs := workflow.WithActivityOptions(ctx, ao)
	var disks []hestia_compute_resources.OrgResourceDisks
	err = workflow.ExecuteActivity(getDisksAtCloudCtxNs, c.CreateSetupTopologyActivities.SelectDiskResourcesAtCloudCtxNs, params.OrgUser.OrgID, params.Kns.CloudCtxNs).Get(getDisksAtCloudCtxNs, &disks)
	if err != nil {
		log.Error("Failed to get disk resources at cloud ctx ns", "Error", err)
		return err
	}
	if len(disks) > 0 {
		resourceIDs := make([]int, len(disks))
		for i, disk := range disks {
			resourceIDs[i] = disk.OrgResources.OrgResourceID
		}
		endDiskServiceCtx := workflow.WithActivityOptions(ctx, ao)
		req := base_deploy_params.DestroyResourcesRequest{
			Ou:             params.OrgUser,
			OrgResourceIDs: resourceIDs,
		}
		err = workflow.ExecuteActivity(endDiskServiceCtx, c.CreateSetupTopologyActivities.EndResourceService, req).Get(endDiskServiceCtx, nil)
		if err != nil {
			log.Error("Failed to update org_resources to end disk service", "Error", err)
			return err
		}
	}
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
	}
	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)
	childWorkflowFuture := workflow.ExecuteChildWorkflow(ctx, "DestroyDeployedTopologyWorkflow", params)
	var childWE workflow.Execution
	if err = childWorkflowFuture.GetChildWorkflowExecution().Get(ctx, &childWE); err != nil {
		return err
	}
	return nil
}
