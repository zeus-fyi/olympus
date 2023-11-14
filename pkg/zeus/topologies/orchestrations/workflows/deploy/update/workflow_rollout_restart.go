package deploy_workflow_cluster_updates

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/definitions/state"
	"github.com/zeus-fyi/zeus/zeus/workload_config_drivers/topology_workloads"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// DeployRolloutRestartWorkflow is a workflow that restarts all workloads in a topology at a single context-cloud-namespace
func (t *FleetUpgradeWorkflow) DeployRolloutRestartWorkflow(ctx workflow.Context, wfID string, ou org_users.OrgUser, cloudCtxNsId int, cloudCtxNs zeus_common_types.CloudCtxNs) error {
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

	var topologyView read_topology_deployment_status.ReadDeploymentStatusesGroup
	getTopologiesAtCloudCtxNsCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(getTopologiesAtCloudCtxNsCtx, t.TopologyUpdateActivity.GetClusterTopologyAtCloudCtxNs, ou, cloudCtxNsId).Get(getTopologiesAtCloudCtxNsCtx, &topologyView)
	if err != nil {
		logger.Error("Failed to get ClusterTopologyAtCloudCtxNs", "Error", err)
		return err
	}

	for _, topology := range topologyView.Slice {
		var infraConfig *topology_workloads.TopologyBaseInfraWorkload
		deployStatusCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(deployStatusCtx, "GetTopologyInfraConfig", ou, topology.TopologyID).Get(deployStatusCtx, &infraConfig)
		if err != nil {
			logger.Error("Failed to get topology infra config", "Error", err)
			return err
		}
		if infraConfig == nil {
			continue
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

	return nil
}
