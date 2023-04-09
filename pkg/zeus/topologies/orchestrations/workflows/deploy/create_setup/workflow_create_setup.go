package deploy_workflow_cluster_setup

import (
	"time"

	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	deploy_topology_activities_create_setup "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/cluster_setup"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"go.temporal.io/sdk/workflow"
)

type ClusterSetupWorkflow struct {
	temporal_base.Workflow
	deploy_topology_activities_create_setup.CreateSetupTopologyActivities
}

const defaultTimeout = 10 * time.Minute

func NewDeployCreateSetupTopologyWorkflow() ClusterSetupWorkflow {
	deployWf := ClusterSetupWorkflow{
		Workflow:                      temporal_base.Workflow{},
		CreateSetupTopologyActivities: deploy_topology_activities_create_setup.CreateSetupTopologyActivities{},
	}
	return deployWf
}

func (c *ClusterSetupWorkflow) GetWorkflows() []interface{} {
	return []interface{}{c.DeployClusterSetupWorkflow}
}

func (c *ClusterSetupWorkflow) DeployClusterSetupWorkflow(ctx workflow.Context, params base_deploy_params.TopologyWorkflowRequest) error {
	log := workflow.GetLogger(ctx)

	c.CreateSetupTopologyActivities.TopologyWorkflowRequest = params
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}

	// TODO add billing email step

	// TODO params
	authCloudCtxNsCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(authCloudCtxNsCtx, c.CreateSetupTopologyActivities.AddAuthCtxNsOrg).Get(authCloudCtxNsCtx, nil)
	if err != nil {
		log.Error("Failed to authorize auth ns to org account", "Error", err)
		return err
	}

	// TODO params
	nodePoolRequestStatusCtxKns := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(nodePoolRequestStatusCtxKns, c.CreateSetupTopologyActivities.MakeNodePoolRequest).Get(nodePoolRequestStatusCtxKns, nil)
	if err != nil {
		log.Error("Failed to complete node pool request", "Error", err)
		return err
	}
	// TODO params
	nodePoolOrgResourcesCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(nodePoolRequestStatusCtxKns, c.CreateSetupTopologyActivities.AddNodePoolToOrgResources).Get(nodePoolOrgResourcesCtx, nil)
	if err != nil {
		log.Error("Failed to add node resources to org account", "Error", err)
		return err
	}

	// TODO params
	diskOrgResourcesCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(diskOrgResourcesCtx, c.CreateSetupTopologyActivities.AddDiskResourcesToOrg).Get(diskOrgResourcesCtx, nil)
	if err != nil {
		log.Error("Failed to add disk resources to org account", "Error", err)
		return err
	}

	//// TODO params
	//domainRequestCtx := workflow.WithActivityOptions(ctx, ao)
	//err = workflow.ExecuteActivity(domainRequestCtx, c.CreateSetupTopologyActivities.AddDomainRecord).Get(domainRequestCtx, nil)
	//if err != nil {
	//	log.Error("Failed to add subdomain resources to org account", "Error", err)
	//	return err
	//}

	emailStatusCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(emailStatusCtx, c.CreateSetupTopologyActivities.SendEmailNotification).Get(emailStatusCtx, nil)
	if err != nil {
		log.Error("Failed to send email notification", "Error", err)
		return err
	}
	return err
}
