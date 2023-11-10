package deploy_workflow_cluster_setup

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	do_types "github.com/zeus-fyi/olympus/pkg/hestia/digitalocean/types"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	kronos_helix "github.com/zeus-fyi/olympus/pkg/kronos/helix"
	deploy_topology_activities_create_setup "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/cluster_setup"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	zeus_templates "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create/templates"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type ClusterSetupWorkflows struct {
	temporal_base.Workflow
	deploy_topology_activities_create_setup.CreateSetupTopologyActivities
}

const defaultTimeout = 120 * time.Minute

func NewDeployCreateSetupTopologyWorkflow() ClusterSetupWorkflows {
	deployWf := ClusterSetupWorkflows{
		Workflow: temporal_base.Workflow{},
		CreateSetupTopologyActivities: deploy_topology_activities_create_setup.CreateSetupTopologyActivities{
			KronosActivities: kronos_helix.NewKronosActivities(),
		},
	}
	return deployWf
}

func (c *ClusterSetupWorkflows) GetDeployClusterSetupWorkflow() interface{} {
	return c.DeployClusterSetupWorkflow
}

func (c *ClusterSetupWorkflows) GetWorkflows() []interface{} {
	return []interface{}{c.DeployClusterSetupWorkflow}
}

// TODO, make app taints optional

func (c *ClusterSetupWorkflows) DeployClusterSetupWorkflow(ctx workflow.Context, wfID string, params base_deploy_params.ClusterSetupRequest) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	oj := artemis_orchestrations.NewInternalActiveTemporalOrchestrationJobTemplate(wfID, "ClusterSetupWorkflows", "DeployClusterSetupWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	aerr := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if aerr != nil {
		logger.Error("Failed to upsert assignment", "Error", aerr)
		return aerr
	}
	// TODO add billing email step
	if params.NodesQuantity > 0 {
		switch params.CloudCtxNs.CloudProvider {
		case "do":
			nodePoolRequestStatusCtxKns := workflow.WithActivityOptions(ctx, ao)
			var nodePoolRequestStatus do_types.DigitalOceanNodePoolRequestStatus
			err := workflow.ExecuteActivity(nodePoolRequestStatusCtxKns, c.CreateSetupTopologyActivities.MakeNodePoolRequest, params).Get(nodePoolRequestStatusCtxKns, &nodePoolRequestStatus)
			if err != nil {
				logger.Error("Failed to complete node pool request for do", "Error", err)
				return err
			}
			nodePoolOrgResourcesCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(nodePoolOrgResourcesCtx, c.CreateSetupTopologyActivities.AddNodePoolToOrgResources, params, nodePoolRequestStatus).Get(nodePoolOrgResourcesCtx, nil)
			if err != nil {
				logger.Error("Failed to add node resources to org account for do", "Error", err)
				return err
			}
		case "gcp":
			nodePoolRequestStatusCtxKns := workflow.WithActivityOptions(ctx, ao)
			var nodePoolRequestStatus do_types.DigitalOceanNodePoolRequestStatus
			err := workflow.ExecuteActivity(nodePoolRequestStatusCtxKns, c.CreateSetupTopologyActivities.GkeMakeNodePoolRequest, params).Get(nodePoolRequestStatusCtxKns, &nodePoolRequestStatus)
			if err != nil {
				logger.Error("Failed to complete node pool request for gke", "Error", err)
				return err
			}
			nodePoolOrgResourcesCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(nodePoolOrgResourcesCtx, c.CreateSetupTopologyActivities.GkeAddNodePoolToOrgResources, params, nodePoolRequestStatus).Get(nodePoolOrgResourcesCtx, nil)
			if err != nil {
				logger.Error("Failed to add node resources to org account for gke", "Error", err)
				return err
			}
		case "aws":
			nodePoolRequestStatusCtxKns := workflow.WithActivityOptions(ctx, ao)
			var nodePoolRequestStatus do_types.DigitalOceanNodePoolRequestStatus
			err := workflow.ExecuteActivity(nodePoolRequestStatusCtxKns, c.CreateSetupTopologyActivities.EksMakeNodePoolRequest, params).Get(nodePoolRequestStatusCtxKns, &nodePoolRequestStatus)
			if err != nil {
				logger.Error("Failed to complete node pool request for eks", "Error", err)
				return err
			}
			nodePoolOrgResourcesCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(nodePoolOrgResourcesCtx, c.CreateSetupTopologyActivities.EksAddNodePoolToOrgResources, params, nodePoolRequestStatus).Get(nodePoolOrgResourcesCtx, nil)
			if err != nil {
				logger.Error("Failed to add node resources to org account for eks", "Error", err)
				return err
			}
		case "ovh":
			nodePoolRequestStatusCtxKns := workflow.WithActivityOptions(ctx, ao)
			var nodePoolRequestStatus do_types.DigitalOceanNodePoolRequestStatus
			err := workflow.ExecuteActivity(nodePoolRequestStatusCtxKns, c.CreateSetupTopologyActivities.OvhMakeNodePoolRequest, params).Get(nodePoolRequestStatusCtxKns, &nodePoolRequestStatus)
			if err != nil {
				logger.Error("Failed to complete node pool request for ovh", "Error", err)
				return err
			}
			nodePoolOrgResourcesCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(nodePoolOrgResourcesCtx, c.CreateSetupTopologyActivities.OvhAddNodePoolToOrgResources, params, nodePoolRequestStatus).Get(nodePoolOrgResourcesCtx, nil)
			if err != nil {
				logger.Error("Failed to add node resources to org account for ovh", "Error", err)
				return err
			}
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
		logger.Error("Failed to authorize auth ns to org account", "Error", err)
		return err
	}

	// TODO needs to add option for gcp, and aws
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
			logger.Error("Failed to add disk resources to org account", "Error", err)
			return err
		}
	}

	//emailStatusCtx := workflow.WithActivityOptions(ctx, ao)
	//err = workflow.ExecuteActivity(emailStatusCtx, c.CreateSetupTopologyActivities.SendEmailNotification, params).Get(emailStatusCtx, nil)
	//if err != nil {
	//	logger.Error("Failed to send email notification", "Error", err)
	//	return err
	//}
	domainRequestCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(domainRequestCtx, c.CreateSetupTopologyActivities.AddDomainRecord, params.CloudCtxNs).Get(domainRequestCtx, nil)
	if err != nil {
		logger.Error("Failed to add subdomain resources to org account", "Error", err)
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

	var ct read_topology.ClusterTopology
	getClusterTopologyIdsCtx := workflow.WithActivityOptions(ctx, aoDeploy)
	err = workflow.ExecuteActivity(getClusterTopologyIdsCtx, c.CreateSetupTopologyActivities.GetClusterTopologyIds, params.Cluster.ClusterName, sbNames, params.Ou).Get(getClusterTopologyIdsCtx, &ct)
	if err != nil {
		logger.Error("Failed to deploy cluster", "Error", err)
		return err
	}
	params.AppTaint = true
	if len(ct.ClusterClassName) <= 0 {
		params.AppTaint = false
	}
	wfParams := base_deploy_params.ClusterTopologyWorkflowRequest{
		ClusterClassName:          ct.ClusterClassName,
		TopologyIDs:               ct.GetTopologyIDs(),
		CloudCtxNS:                params.CloudCtxNs,
		OrgUser:                   params.Ou,
		AppTaint:                  params.AppTaint,
		RequestChoreographySecret: ct.CheckForChoreographyOption(),
	}
	deployChildWorkflowOptions := workflow.ChildWorkflowOptions{
		ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
	}
	clusterDeployCtx := workflow.WithChildOptions(ctx, deployChildWorkflowOptions)
	deployChildWorkflowFuture := workflow.ExecuteChildWorkflow(clusterDeployCtx, "DeployClusterTopologyWorkflow", wfID, wfParams)
	var deployChildWfExec workflow.Execution
	if err = deployChildWorkflowFuture.GetChildWorkflowExecution().Get(ctx, &deployChildWfExec); err != nil {
		logger.Error("Failed to get child deployment workflow execution", "Error", err)
		return err
	}

	childWorkflowOptions := workflow.ChildWorkflowOptions{
		ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
	}
	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)
	ftDestroy := base_deploy_params.DestroyClusterSetupRequest{
		ClusterSetupRequest: base_deploy_params.ClusterSetupRequest{
			FreeTrial:     params.FreeTrial,
			Ou:            params.Ou,
			CloudCtxNs:    params.CloudCtxNs,
			Nodes:         params.Nodes,
			NodesQuantity: params.NodesQuantity,
			Disks:         params.Disks,
			Cluster: zeus_templates.Cluster{
				ClusterName:     params.Cluster.ClusterName,
				ComponentBases:  nil,
				IngressSettings: zeus_templates.Ingress{},
				IngressPaths:    nil,
			},
			AppTaint: wfParams.AppTaint,
		},
	}
	childWorkflowFuture := workflow.ExecuteChildWorkflow(ctx, "DestroyClusterSetupWorkflowFreeTrial", wfID, ftDestroy)
	var childWE workflow.Execution
	if err = childWorkflowFuture.GetChildWorkflowExecution().Get(ctx, &childWE); err != nil {
		logger.Error("Failed to get child workflow execution", "Error", err)
		return err
	}
	finishedCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
	if err != nil {
		logger.Error("Failed to update and mark orchestration inactive", "Error", err)
		return err
	}
	return nil
}
