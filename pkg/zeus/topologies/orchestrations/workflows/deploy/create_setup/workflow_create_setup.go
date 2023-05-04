package deploy_workflow_cluster_setup

import (
	"time"

	do_types "github.com/zeus-fyi/olympus/pkg/hestia/digitalocean/types"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	deploy_topology_activities_create_setup "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/cluster_setup"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type ClusterSetupWorkflow struct {
	temporal_base.Workflow
	deploy_topology_activities_create_setup.CreateSetupTopologyActivities
}

const defaultTimeout = 120 * time.Minute

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
	switch params.CloudCtxNs.CloudProvider {
	case "do":
		nodePoolRequestStatusCtxKns := workflow.WithActivityOptions(ctx, ao)
		var nodePoolRequestStatus do_types.DigitalOceanNodePoolRequestStatus
		err := workflow.ExecuteActivity(nodePoolRequestStatusCtxKns, c.CreateSetupTopologyActivities.MakeNodePoolRequest, params).Get(nodePoolRequestStatusCtxKns, &nodePoolRequestStatus)
		if err != nil {
			log.Error("Failed to complete node pool request", "Error", err)
			return err
		}
		nodePoolOrgResourcesCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(nodePoolOrgResourcesCtx, c.CreateSetupTopologyActivities.AddNodePoolToOrgResources, params, nodePoolRequestStatus).Get(nodePoolOrgResourcesCtx, nil)
		if err != nil {
			log.Error("Failed to add node resources to org account", "Error", err)
			return err
		}
	case "gcp":
		nodePoolRequestStatusCtxKns := workflow.WithActivityOptions(ctx, ao)
		var nodePoolRequestStatus do_types.DigitalOceanNodePoolRequestStatus
		err := workflow.ExecuteActivity(nodePoolRequestStatusCtxKns, c.CreateSetupTopologyActivities.GkeMakeNodePoolRequest, params).Get(nodePoolRequestStatusCtxKns, &nodePoolRequestStatus)
		if err != nil {
			log.Error("Failed to complete node pool request", "Error", err)
			return err
		}
		nodePoolOrgResourcesCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(nodePoolOrgResourcesCtx, c.CreateSetupTopologyActivities.GkeAddNodePoolToOrgResources, params, nodePoolRequestStatus).Get(nodePoolOrgResourcesCtx, nil)
		if err != nil {
			log.Error("Failed to add node resources to org account", "Error", err)
			return err
		}
	}
	authCloudCtxNsCtxOptions := ao
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second * 60,
		BackoffCoefficient: 1.0,
		MaximumInterval:    time.Second * 60,
	}
	authCloudCtxNsCtxOptions.RetryPolicy = retryPolicy
	authCloudCtxNsCtx := workflow.WithActivityOptions(ctx, authCloudCtxNsCtxOptions)
	err := workflow.ExecuteActivity(authCloudCtxNsCtx, c.CreateSetupTopologyActivities.AddAuthCtxNsOrg, params).Get(authCloudCtxNsCtx, nil)
	if err != nil {
		log.Error("Failed to authorize auth ns to org account", "Error", err)
		return err
	}

	// TODO needs to add option for gcp
	for _, disk := range params.Disks {
		if disk.DiskSize == 0 && disk.DiskUnits == "" {
			continue
		}
		diskActivityOptions := ao
		retryPolicy = &temporal.RetryPolicy{
			InitialInterval:    time.Second * 60,
			BackoffCoefficient: 1.0,
			MaximumInterval:    time.Second * 60,
		}
		diskActivityOptions.RetryPolicy = retryPolicy
		diskOrgResourcesCtx := workflow.WithActivityOptions(ctx, diskActivityOptions)
		err = workflow.ExecuteActivity(diskOrgResourcesCtx, c.CreateSetupTopologyActivities.AddDiskResourcesToOrg, params, disk).Get(diskOrgResourcesCtx, nil)
		if err != nil {
			log.Error("Failed to add disk resources to org account", "Error", err)
			return err
		}
	}
	//emailStatusCtx := workflow.WithActivityOptions(ctx, ao)
	//err = workflow.ExecuteActivity(emailStatusCtx, c.CreateSetupTopologyActivities.SendEmailNotification, params).Get(emailStatusCtx, nil)
	//if err != nil {
	//	log.Error("Failed to send email notification", "Error", err)
	//	return err
	//}
	domainRequestCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(domainRequestCtx, c.CreateSetupTopologyActivities.AddDomainRecord, params.Namespace).Get(domainRequestCtx, nil)
	if err != nil {
		log.Error("Failed to add subdomain resources to org account", "Error", err)
		return err
	}
	deployRetryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second * 15,
		BackoffCoefficient: 2,
	}
	aoDeploy := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 15,
		RetryPolicy:         deployRetryPolicy,
	}
	var sbNames []string
	for _, cb := range params.Cluster.ComponentBases {
		for sbName, _ := range cb {
			sbNames = append(sbNames, sbName)
		}
	}
	clusterDeployCtx := workflow.WithActivityOptions(ctx, aoDeploy)
	err = workflow.ExecuteActivity(clusterDeployCtx, c.CreateSetupTopologyActivities.DeployClusterTopologyFromUI, params.Cluster.ClusterName, sbNames, params.CloudCtxNs, params.Ou).Get(clusterDeployCtx, nil)
	if err != nil {
		log.Error("Failed to deploy cluster", "Error", err)
		return err
	}
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
	}
	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)
	ftDestroy := base_deploy_params.DestroyClusterSetupRequest{
		ClusterSetupRequest: params,
	}
	childWorkflowFuture := workflow.ExecuteChildWorkflow(ctx, "DestroyClusterSetupWorkflowFreeTrial", ftDestroy)
	var childWE workflow.Execution
	if err = childWorkflowFuture.GetChildWorkflowExecution().Get(ctx, &childWE); err != nil {
		log.Error("Failed to get child workflow execution", "Error", err)
		return err
	}
	return nil
}
