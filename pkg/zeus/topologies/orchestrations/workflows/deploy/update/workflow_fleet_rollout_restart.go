package deploy_workflow_cluster_updates

import (
	"time"

	read_topologies "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies"
	read_topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/definitions/state"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_resp_types/topology_workloads"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	deployment  = "deployment"
	statefulset = "statefulset"
)

func (t *FleetUpgradeWorkflow) DeployRolloutRestartFleetWorkflow(ctx workflow.Context, wfID string, params base_deploy_params.FleetRolloutRestartWorkflowRequest) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 60,
			BackoffCoefficient: 1.0,
			MaximumInterval:    time.Second * 60,
			MaximumAttempts:    1000,
		},
	}
	var clusterToUpgrade []read_topologies.ClusterAppView
	workerCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(workerCtx, t.TopologyUpdateActivity.GetClustersToUpdate, params.OrgUser, params.ClusterName).Get(workerCtx, &clusterToUpgrade)
	if err != nil {
		logger.Error("Failed to get clusters to update", "Error", err)
		return err
	}
	for _, clusterInfo := range clusterToUpgrade {
		var topologyView read_topology_deployment_status.ReadDeploymentStatusesGroup
		getTopologiesAtCloudCtxNsCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(getTopologiesAtCloudCtxNsCtx, t.TopologyUpdateActivity.GetClusterTopologyAtCloudCtxNs, params.OrgUser, clusterInfo.CloudCtxNsID).Get(getTopologiesAtCloudCtxNsCtx, &topologyView)
		if err != nil {
			logger.Error("Failed to get ClusterTopologyAtCloudCtxNs", "Error", err)
			return err
		}

		for _, topology := range topologyView.Slice {
			var infraConfig *topology_workloads.TopologyBaseInfraWorkload
			deployStatusCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(deployStatusCtx, "GetTopologyInfraConfig", params.OrgUser, topology.TopologyID).Get(deployStatusCtx, &infraConfig)
			if err != nil {
				logger.Error("Failed to get topology infra config", "Error", err)
				return err
			}
			if infraConfig == nil {
				continue
			}
			cloudCtxNs := zeus_common_types.CloudCtxNs{
				CloudProvider: clusterInfo.CloudProvider,
				Region:        clusterInfo.Region,
				Context:       clusterInfo.Context,
				Namespace:     clusterInfo.Namespace,
				Alias:         clusterInfo.NamespaceAlias,
			}
			if infraConfig.Deployment != nil && infraConfig.Deployment.Name != "" {
				restartDepCtx := workflow.WithActivityOptions(ctx, ao)
				err = workflow.ExecuteActivity(restartDepCtx, t.RestartWorkload, cloudCtxNs, infraConfig.Deployment.Name, deployment).Get(restartDepCtx, nil)
				if err != nil {
					logger.Error("Failed to get topology infra config", "Error", err)
					return err
				}
			}

			if infraConfig.StatefulSet != nil && infraConfig.StatefulSet.Name != "" {
				restartStsCtx := workflow.WithActivityOptions(ctx, ao)
				err = workflow.ExecuteActivity(restartStsCtx, t.RestartWorkload, cloudCtxNs, infraConfig.StatefulSet.Name, statefulset).Get(restartStsCtx, nil)
				if err != nil {
					logger.Error("Failed to get topology infra config", "Error", err)
					return err
				}
			}
		}

	}
	return nil
}
