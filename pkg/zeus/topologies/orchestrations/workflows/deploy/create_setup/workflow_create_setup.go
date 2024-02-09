package deploy_workflow_cluster_setup

import (
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/authorized_clusters"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	do_types "github.com/zeus-fyi/olympus/pkg/hestia/digitalocean/types"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	kronos_helix "github.com/zeus-fyi/olympus/pkg/kronos/helix"
	deploy_topology_activities_create_setup "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/cluster_setup"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"github.com/zeus-fyi/zeus/zeus/workload_config_drivers/topology_workloads"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
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

	var authCfg *authorized_clusters.K8sClusterConfig
	clusterAuthCtxKns := workflow.WithActivityOptions(ctx, ao)
	cerr := workflow.ExecuteActivity(clusterAuthCtxKns, c.CreateSetupTopologyActivities.GetClusterAuthCtx, params.Ou, params.CloudCtxNs).Get(clusterAuthCtxKns, &authCfg)
	if cerr != nil {
		logger.Error("Failed to get cluster auth ctx", "Error", cerr)
		return cerr
	}

	isPublic := false
	if authCfg != nil {
		tmp := params
		tmp.CloudCtxNs = authCfg.CloudCtxNs
		params = tmp
		isPublic = authCfg.IsPublic
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
			switch isPublic {
			case true:
				err := workflow.ExecuteActivity(nodePoolRequestStatusCtxKns, c.CreateSetupTopologyActivities.EksMakeNodePoolRequest, params).Get(nodePoolRequestStatusCtxKns, &nodePoolRequestStatus)
				if err != nil {
					logger.Error("Failed to complete node pool request for eks", "Error", err)
					return err
				}
				if len(nodePoolRequestStatus.ClusterID) == 0 || len(nodePoolRequestStatus.NodePoolID) == 0 {
					err = fmt.Errorf("failed to get cluster id")
					logger.Error("Failed to get cluster id", "Error", err)
					return err
				}

				nodePoolOrgResourcesCtx := workflow.WithActivityOptions(ctx, ao)
				err = workflow.ExecuteActivity(nodePoolOrgResourcesCtx, c.CreateSetupTopologyActivities.EksAddNodePoolToOrgResources, params, nodePoolRequestStatus).Get(nodePoolOrgResourcesCtx, nil)
				if err != nil {
					logger.Error("Failed to add node resources to org account for eks", "Error", err)
					return err
				}
			case false:
				err := workflow.ExecuteActivity(nodePoolRequestStatusCtxKns, c.CreateSetupTopologyActivities.PrivateEksMakeNodePoolRequest, params).Get(nodePoolRequestStatusCtxKns, &nodePoolRequestStatus)
				if err != nil {
					logger.Error("Failed to complete private node pool request for eks", "Error", err)
					return err
				}
				if len(nodePoolRequestStatus.ClusterID) == 0 || len(nodePoolRequestStatus.NodePoolID) == 0 || nodePoolRequestStatus.ExtClusterCfgID <= 0 {
					err = fmt.Errorf("failed to get cluster id")
					logger.Error("Failed to get cluster id", "Error", err)
					return err
				}
				nodePoolOrgResourcesCtx := workflow.WithActivityOptions(ctx, ao)
				err = workflow.ExecuteActivity(nodePoolOrgResourcesCtx, c.CreateSetupTopologyActivities.EksAddNodePoolToOrgResources, params, nodePoolRequestStatus).Get(nodePoolOrgResourcesCtx, nil)
				if err != nil {
					logger.Error("Failed to add node resources to org account for eks", "Error", err)
					return err
				}
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
		CloudCtxNs:                params.CloudCtxNs,
		OrgUser:                   params.Ou,
		AppTaint:                  params.AppTaint,
		RequestChoreographySecret: ct.CheckForChoreographyOption(),
	}

	for _, topID := range wfParams.TopologyIDs {
		var infraConfig *topology_workloads.TopologyBaseInfraWorkload
		getTopInfraCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(getTopInfraCtx, "GetTopologyInfraConfig", params.Ou, topID).Get(getTopInfraCtx, &infraConfig)
		if err != nil {
			logger.Error("Failed to get topology infra config", "Error", err)
			return err
		}
		if infraConfig == nil {
			logger.Info("Topology infra config is nil", "TopologyID", topID)
			continue
		}
		if topID == 0 {
			logger.Info("Topology id is empty", "ClusterClassName", wfParams.ClusterClassName)
			continue
		}
		desWf := base_deploy_params.TopologyWorkflowRequest{
			TopologyDeployRequest: zeus_req_types.TopologyDeployRequest{
				TopologyID:                topID,
				CloudCtxNs:                params.CloudCtxNs,
				TopologyBaseInfraWorkload: *infraConfig,
			},
			OrgUser: params.Ou,
		}
		childWorkflowOptions := workflow.ChildWorkflowOptions{
			ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
			RetryPolicy:       deployRetryPolicy,
		}
		wfChildCtx := workflow.WithChildOptions(ctx, childWorkflowOptions)
		childWorkflowFuture := workflow.ExecuteChildWorkflow(wfChildCtx, "DeployTopologyWorkflow", wfID, desWf)
		var childWE workflow.Execution
		if err = childWorkflowFuture.GetChildWorkflowExecution().Get(wfChildCtx, &childWE); err != nil {
			logger.Error("Failed to get child workflow execution", "Error", err)
			return err
		}
	}
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		RetryPolicy:       retryPolicy,
		ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
	}
	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)
	childWorkflowFuture := workflow.ExecuteChildWorkflow(ctx, "DestroyClusterSetupWorkflowFreeTrial", wfID, params, wfParams)
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
