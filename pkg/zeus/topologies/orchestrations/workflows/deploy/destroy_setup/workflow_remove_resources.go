package deploy_workflow_destroy_setup

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/authorized_clusters"
	do_types "github.com/zeus-fyi/olympus/pkg/hestia/digitalocean/types"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	kronos_helix "github.com/zeus-fyi/olympus/pkg/kronos/helix"
	deploy_topology_activities_create_setup "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/cluster_setup"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"go.temporal.io/sdk/workflow"
)

type DestroyResourcesWorkflows struct {
	temporal_base.Workflow
	deploy_topology_activities_create_setup.CreateSetupTopologyActivities
}

func NewDestroyResourcesWorkflow() DestroyResourcesWorkflows {
	deployWf := DestroyResourcesWorkflows{
		Workflow: temporal_base.Workflow{},
		CreateSetupTopologyActivities: deploy_topology_activities_create_setup.CreateSetupTopologyActivities{
			KronosActivities: kronos_helix.NewKronosActivities(),
		},
	}
	return deployWf
}

func (c *DestroyResourcesWorkflows) GetWorkflow() interface{} {
	return c.DestroyClusterResourcesWorkflow
}

func (c *DestroyResourcesWorkflows) GetWorkflows() []interface{} {
	return []interface{}{c.DestroyClusterResourcesWorkflow}
}

func (c *DestroyResourcesWorkflows) DestroyClusterResourcesWorkflow(ctx workflow.Context, wfID string, params base_deploy_params.DestroyResourcesRequest) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	oj := artemis_orchestrations.NewInternalActiveTemporalOrchestrationJobTemplate(wfID, "DestroyResourcesWorkflows", "DestroyClusterResourcesWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	aerr := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if aerr != nil {
		logger.Error("Failed to upsert assignment", "Error", aerr)
		return aerr
	}
	selectNodesCtx := workflow.WithActivityOptions(ctx, ao)
	var nodes []do_types.DigitalOceanNodePoolRequestStatus
	err := workflow.ExecuteActivity(selectNodesCtx, c.CreateSetupTopologyActivities.SelectNodeResources, params).Get(selectNodesCtx, &nodes)
	if err != nil {
		logger.Error("Failed to select org resource nodes", "Error", err)
		return err
	}
	var gkeNodes []do_types.DigitalOceanNodePoolRequestStatus
	err = workflow.ExecuteActivity(selectNodesCtx, c.CreateSetupTopologyActivities.SelectGkeNodeResources, params).Get(selectNodesCtx, &gkeNodes)
	if err != nil {
		logger.Error("Failed to select org resource nodes", "Error", err)
		return err
	}
	var eksNodes []do_types.DigitalOceanNodePoolRequestStatus
	err = workflow.ExecuteActivity(selectNodesCtx, c.CreateSetupTopologyActivities.SelectEksNodeResources, params).Get(selectNodesCtx, &eksNodes)
	if err != nil {
		logger.Error("Failed to select org resource nodes", "Error", err)
		return err
	}
	var ovhNodes []do_types.DigitalOceanNodePoolRequestStatus
	err = workflow.ExecuteActivity(selectNodesCtx, c.CreateSetupTopologyActivities.SelectOvhNodeResources, params).Get(selectNodesCtx, &ovhNodes)
	if err != nil {
		logger.Error("Failed to select org resource nodes for ovh", "Error", err)
		return err
	}
	if len(nodes) == 0 && len(gkeNodes) == 0 && len(eksNodes) == 0 && len(ovhNodes) == 0 {
		logger.Info("No node resources found to destroy or they were free trial nodes that will be deleted automatically")
		return nil
	}
	for _, node := range nodes {
		destroyNodePoolOrgResourcesCtx := workflow.WithActivityOptions(ctx, ao)
		logger.Info("Destroying node pool org resources", "NodePoolRequestStatus", node)
		err = workflow.ExecuteActivity(destroyNodePoolOrgResourcesCtx, c.CreateSetupTopologyActivities.RemoveNodePoolRequest, node).Get(destroyNodePoolOrgResourcesCtx, nil)
		if err != nil {
			logger.Error("Failed to remove node resources for account", "Error", err)
			return err
		}
	}
	for _, node := range gkeNodes {
		destroyNodePoolOrgResourcesCtx := workflow.WithActivityOptions(ctx, ao)
		logger.Info("Destroying node pool org resources", "GkeRemoveNodePoolRequest", node)
		err = workflow.ExecuteActivity(destroyNodePoolOrgResourcesCtx, c.CreateSetupTopologyActivities.GkeRemoveNodePoolRequest, node).Get(destroyNodePoolOrgResourcesCtx, nil)
		if err != nil {
			logger.Error("Failed to remove node resources for account", "Error", err)
			return err
		}
	}
	for _, node := range eksNodes {
		if node.ExtClusterCfgID > 0 {
			eksDestroyNodePoolOrgResourcesCtx := workflow.WithActivityOptions(ctx, ao)
			logger.Info("Destroying node pool org resources", "PrivateEksRemoveNodePoolRequest", node)

			var authCfg *authorized_clusters.K8sClusterConfig
			clusterAuthCtxKns := workflow.WithActivityOptions(ctx, ao)
			cerr := workflow.ExecuteActivity(clusterAuthCtxKns, c.CreateSetupTopologyActivities.GetClusterAuthCtxFromID, params.Ou, fmt.Sprintf("%d", node.ExtClusterCfgID)).Get(clusterAuthCtxKns, &authCfg)
			if cerr != nil {
				logger.Error("Failed to get cluster auth ctx", "Error", cerr)
				return cerr
			}

			if authCfg == nil || authCfg.Context == "" || authCfg.Region == "" {
				cerr = fmt.Errorf("failed to get cluster auth ctx")
				logger.Error("Failed to get cluster auth ctx", "Error", cerr)
				return cerr
			}

			// TOOD, lookup cloudCtxNs
			err = workflow.ExecuteActivity(eksDestroyNodePoolOrgResourcesCtx, c.CreateSetupTopologyActivities.PrivateEksRemoveNodePoolRequest, params.Ou, authCfg.CloudCtxNs, node).Get(eksDestroyNodePoolOrgResourcesCtx, nil)
			if err != nil {
				logger.Error("Failed to remove eks node resources for account", "Error", err)
				return err
			}
		} else {
			destroyNodePoolOrgResourcesCtx := workflow.WithActivityOptions(ctx, ao)
			logger.Info("Destroying node pool org resources", "EksRemoveNodePoolRequest", node)
			err = workflow.ExecuteActivity(destroyNodePoolOrgResourcesCtx, c.CreateSetupTopologyActivities.EksRemoveNodePoolRequest, node).Get(destroyNodePoolOrgResourcesCtx, nil)
			if err != nil {
				logger.Error("Failed to remove node resources for account", "Error", err)
				return err
			}
		}
	}
	for _, node := range ovhNodes {
		destroyNodePoolOrgResourcesCtx := workflow.WithActivityOptions(ctx, ao)
		logger.Info("Destroying node pool org resources", "OvhRemoveNodePoolRequest", node)
		err = workflow.ExecuteActivity(destroyNodePoolOrgResourcesCtx, c.CreateSetupTopologyActivities.OvhRemoveNodePoolRequest, node).Get(destroyNodePoolOrgResourcesCtx, nil)
		if err != nil {
			logger.Error("Failed to remove node resources for account for ovh", "Error", err)
			return err
		}
	}
	// TODO billing somewhere usage update
	endServiceNodesCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(endServiceNodesCtx, c.CreateSetupTopologyActivities.EndResourceService, params).Get(endServiceNodesCtx, &nodes)
	if err != nil {
		logger.Error("Failed to update org_resources to end service", "Error", err)
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
