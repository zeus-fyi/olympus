package deploy_workflow_cluster_setup

import (
	"time"

	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	deploy_topology_activities_create_setup "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/cluster_setup"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	create_or_update_deploy "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/external/create_or_update"
	"go.temporal.io/sdk/workflow"
)

type ClusterSetupWorkflow struct {
	temporal_base.Workflow
	deploy_topology_activities_create_setup.CreateSetupTopologyActivities
}

const defaultTimeout = 60 * time.Minute

func NewDeployCreateSetupTopologyWorkflow() ClusterSetupWorkflow {
	deployWf := ClusterSetupWorkflow{
		Workflow:                      temporal_base.Workflow{},
		CreateSetupTopologyActivities: deploy_topology_activities_create_setup.CreateSetupTopologyActivities{},
	}
	return deployWf
}

func (c *ClusterSetupWorkflow) GetDeployClusterSetupWorkflow() interface{} {
	return c.DeployClusterSetupWorkflow
}

func (c *ClusterSetupWorkflow) GetWorkflows() []interface{} {
	return []interface{}{c.DeployClusterSetupWorkflow}
}

func (c *ClusterSetupWorkflow) DeployClusterSetupWorkflow(ctx workflow.Context, params base_deploy_params.ClusterSetupRequest) error {
	log := workflow.GetLogger(ctx)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}

	// TODO add billing email step
	authCloudCtxNsCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(authCloudCtxNsCtx, c.CreateSetupTopologyActivities.AddAuthCtxNsOrg, params).Get(authCloudCtxNsCtx, nil)
	if err != nil {
		log.Error("Failed to authorize auth ns to org account", "Error", err)
		return err
	}

	nodePoolRequestStatusCtxKns := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(nodePoolRequestStatusCtxKns, c.CreateSetupTopologyActivities.MakeNodePoolRequest, params).Get(nodePoolRequestStatusCtxKns, nil)
	if err != nil {
		log.Error("Failed to complete node pool request", "Error", err)
		return err
	}
	nodePoolOrgResourcesCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(nodePoolRequestStatusCtxKns, c.CreateSetupTopologyActivities.AddNodePoolToOrgResources, params).Get(nodePoolOrgResourcesCtx, nil)
	if err != nil {
		log.Error("Failed to add node resources to org account", "Error", err)
		return err
	}

	diskOrgResourcesCtx := workflow.WithActivityOptions(ctx, ao)
	for _, disk := range params.Disks {
		err = workflow.ExecuteActivity(diskOrgResourcesCtx, c.CreateSetupTopologyActivities.AddDiskResourcesToOrg, params, disk).Get(diskOrgResourcesCtx, nil)
		if err != nil {
			log.Error("Failed to add disk resources to org account", "Error", err)
			return err
		}
	}

	//// TODO params
	//domainRequestCtx := workflow.WithActivityOptions(ctx, ao)
	//err = workflow.ExecuteActivity(domainRequestCtx, c.CreateSetupTopologyActivities.AddDomainRecord, params).Get(domainRequestCtx, nil)
	//if err != nil {
	//	log.Error("Failed to add subdomain resources to org account", "Error", err)
	//	return err
	//}

	//emailStatusCtx := workflow.WithActivityOptions(ctx, ao)
	//err = workflow.ExecuteActivity(emailStatusCtx, c.CreateSetupTopologyActivities.SendEmailNotification, params).Get(emailStatusCtx, nil)
	//if err != nil {
	//	log.Error("Failed to send email notification", "Error", err)
	//	return err
	//}

	var sbNames []string
	for _, cb := range params.Cluster.ComponentBases {
		for sbName, _ := range cb {
			sbNames = append(sbNames, sbName)
		}
	}
	cdRequest := create_or_update_deploy.TopologyClusterDeployRequest{
		ClusterClassName:    params.Cluster.ClusterName,
		SkeletonBaseOptions: sbNames,
		CloudCtxNs:          params.CloudCtxNs,
	}

	clusterDeployCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(clusterDeployCtx, c.CreateSetupTopologyActivities.DeployClusterTopology, cdRequest, params.Ou).Get(clusterDeployCtx, nil)
	if err != nil {
		log.Error("Failed to add deploy cluster", "Error", err)
		return err
	}
	return err
}
