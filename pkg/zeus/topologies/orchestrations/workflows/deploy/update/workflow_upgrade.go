package deploy_workflow_cluster_updates

import (
	"time"

	"github.com/google/uuid"
	read_topologies "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies"
	read_topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/definitions/state"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	deploy_topology_update_activities "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/update"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

type FleetUpgradeWorkflow struct {
	temporal_base.Workflow
	deploy_topology_update_activities.TopologyUpdateActivity
}

const defaultTimeout = 60 * time.Minute

func NewDeployFleetUpgradeWorkflow() FleetUpgradeWorkflow {
	deployWf := FleetUpgradeWorkflow{
		Workflow:               temporal_base.Workflow{},
		TopologyUpdateActivity: deploy_topology_update_activities.TopologyUpdateActivity{},
	}
	return deployWf
}

func (t *FleetUpgradeWorkflow) GetWorkflows() []interface{} {
	return []interface{}{t.UpgradeFleetWorkflow}
}

func (t *FleetUpgradeWorkflow) UpgradeFleetWorkflow(ctx workflow.Context, params base_deploy_params.FleetUpgradeWorkflowRequest) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}

	var clusterToUpgrade []read_topologies.ClusterAppView
	workerCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(workerCtx, t.TopologyUpdateActivity.GetClustersToUpdate, params).Get(workerCtx, &clusterToUpgrade)
	if err != nil {
		log.Error("Failed to get clusters to update", "Error", err)
		return err
	}
	for _, clusterInfo := range clusterToUpgrade {
		var topologyView read_topology_deployment_status.ReadDeploymentStatusesGroup
		getTopologiesAtCloudCtxNsCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(getTopologiesAtCloudCtxNsCtx, t.TopologyUpdateActivity.GetClusterTopologyAtCloudCtxNs, params, clusterInfo).Get(getTopologiesAtCloudCtxNsCtx, &topologyView)
		if err != nil {
			log.Error("Failed to get ClusterTopologyAtCloudCtxNs", "Error", err)
			return err
		}
		var clusterReq base_deploy_params.ClusterTopologyWorkflowRequest
		diffTopologiesAtCloudCtxNsCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(diffTopologiesAtCloudCtxNsCtx, t.TopologyUpdateActivity.DiffClusterUpdate, params, clusterInfo, topologyView).Get(diffTopologiesAtCloudCtxNsCtx, &clusterReq)
		if err != nil {
			log.Error("Failed to DiffClusterUpdate", "Error", err)
			return err
		}
		if clusterReq.ClusterClassName == "" {
			continue
		}
		childWorkflowOptions := workflow.ChildWorkflowOptions{
			ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
		}
		childCtx := workflow.WithChildOptions(ctx, childWorkflowOptions)
		childWorkflowFuture := workflow.ExecuteChildWorkflow(childCtx, "DeployClusterTopologyWorkflow", uuid.New().String(), clusterReq)
		var childWE workflow.Execution
		if err = childWorkflowFuture.GetChildWorkflowExecution().Get(ctx, &childWE); err != nil {
			log.Error("Failed to get child workflow execution", "Error", err)
			return err
		}
	}
	return nil
}
