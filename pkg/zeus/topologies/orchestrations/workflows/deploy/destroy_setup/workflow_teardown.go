package deploy_workflow_destroy_setup

import (
	"context"
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/authorized_clusters"
	do_types "github.com/zeus-fyi/olympus/pkg/hestia/digitalocean/types"
	hestia_stripe "github.com/zeus-fyi/olympus/pkg/hestia/stripe"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	deploy_topology_activities_create_setup "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/cluster_setup"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"github.com/zeus-fyi/zeus/zeus/workload_config_drivers/topology_workloads"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

type DestroyClusterSetupWorkflow struct {
	temporal_base.Workflow
	deploy_topology_activities_create_setup.CreateSetupTopologyActivities
}

const defaultTimeout = 120 * time.Minute

func NewDeployDestroyClusterSetupWorkflow() DestroyClusterSetupWorkflow {
	deployWf := DestroyClusterSetupWorkflow{
		Workflow:                      temporal_base.Workflow{},
		CreateSetupTopologyActivities: deploy_topology_activities_create_setup.CreateSetupTopologyActivities{},
	}
	return deployWf
}

func (c *DestroyClusterSetupWorkflow) GetWorkflow() interface{} {
	return c.DestroyClusterSetupWorkflowFreeTrial
}

func (c *DestroyClusterSetupWorkflow) GetWorkflows() []interface{} {
	return []interface{}{c.DestroyClusterSetupWorkflowFreeTrial}
}

func (c *DestroyClusterSetupWorkflow) DestroyClusterSetupWorkflowFreeTrial(ctx workflow.Context, wfID string, params base_deploy_params.ClusterSetupRequest, wfParams base_deploy_params.ClusterTopologyWorkflowRequest) error {
	logger := workflow.GetLogger(ctx)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}

	oj := artemis_orchestrations.NewInternalActiveTemporalOrchestrationJobTemplate(wfID, "DestroyClusterSetupWorkflow", "DestroyClusterSetupWorkflowFreeTrial")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	aerr := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if aerr != nil {
		logger.Error("Failed to upsert assignment", "Error", aerr)
		return aerr
	}

	if params.FreeTrial {
		err := workflow.Sleep(ctx, 60*time.Minute)
		if err != nil {
			logger.Error("Failed to sleep for 1 hour", "Error", err)
			return err
		}
		hestiaCtx := context.Background()
		isBillingSetup, herr := hestia_stripe.DoesUserHaveBillingMethod(hestiaCtx, params.Ou.UserID)
		if herr != nil {
			logger.Error("Failed to check if user has billing method", "Error", herr)
			return herr
		}
		if !isBillingSetup {
			removeSubdomainCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(removeSubdomainCtx, c.CreateSetupTopologyActivities.RemoveDomainRecord, params.CloudCtxNs).Get(removeSubdomainCtx, nil)
			if err != nil {
				logger.Error("Failed to remove domain record", "Error", err)
				return err
			}
			removeAuthCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(removeAuthCtx, c.CreateSetupTopologyActivities.RemoveAuthCtxNsOrg, params.Ou.OrgID, params.CloudCtxNs).Get(removeAuthCtx, nil)
			if err != nil {
				logger.Error("Failed to remove auth ctx ns", "Error", err)
				return err
			}

			selectFreeTrialDoNodesCtx := workflow.WithActivityOptions(ctx, ao)
			var nodes []do_types.DigitalOceanNodePoolRequestStatus
			err = workflow.ExecuteActivity(selectFreeTrialDoNodesCtx, c.CreateSetupTopologyActivities.SelectFreeTrialNodes, params.Ou.OrgID).Get(selectFreeTrialDoNodesCtx, &nodes)
			if err != nil {
				logger.Error("Failed to select digital ocean free trial nodes", "Error", err)
				return err
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
			gkeSelectFreeTrialDoNodesCtx := workflow.WithActivityOptions(ctx, ao)
			var gkeNodes []do_types.DigitalOceanNodePoolRequestStatus
			err = workflow.ExecuteActivity(gkeSelectFreeTrialDoNodesCtx, c.CreateSetupTopologyActivities.GkeSelectFreeTrialNodes, params.Ou.OrgID).Get(gkeSelectFreeTrialDoNodesCtx, &gkeNodes)
			if err != nil {
				logger.Error("Failed to select gke free trial nodes", "Error", err)
				return err
			}
			for _, node := range gkeNodes {
				gkeDestroyNodePoolOrgResourcesCtx := workflow.WithActivityOptions(ctx, ao)
				logger.Info("Destroying node pool org resources", "GkeRemoveNodePoolRequest", node)
				err = workflow.ExecuteActivity(gkeDestroyNodePoolOrgResourcesCtx, c.CreateSetupTopologyActivities.GkeRemoveNodePoolRequest, node).Get(gkeDestroyNodePoolOrgResourcesCtx, nil)
				if err != nil {
					logger.Error("Failed to remove gke node resources for account", "Error", err)
					return err
				}
			}
			ovhSelectFreeTrialDoNodesCtx := workflow.WithActivityOptions(ctx, ao)
			var ovhNodes []do_types.DigitalOceanNodePoolRequestStatus
			err = workflow.ExecuteActivity(ovhSelectFreeTrialDoNodesCtx, c.CreateSetupTopologyActivities.OvhSelectFreeTrialNodes, params.Ou.OrgID).Get(ovhSelectFreeTrialDoNodesCtx, &ovhNodes)
			if err != nil {
				logger.Error("Failed to select ovh free trial nodes", "Error", err)
				return err
			}
			for _, node := range ovhNodes {
				ovhDestroyNodePoolOrgResourcesCtx := workflow.WithActivityOptions(ctx, ao)
				logger.Info("Destroying node pool org resources", "OvhRemoveNodePoolRequest", node)
				err = workflow.ExecuteActivity(ovhDestroyNodePoolOrgResourcesCtx, c.CreateSetupTopologyActivities.OvhRemoveNodePoolRequest, node).Get(ovhDestroyNodePoolOrgResourcesCtx, nil)
				if err != nil {
					logger.Error("Failed to remove ovh node resources for account", "Error", err)
					return err
				}
			}
			eksSelectFreeTrialDoNodesCtx := workflow.WithActivityOptions(ctx, ao)
			var eksNodes []do_types.DigitalOceanNodePoolRequestStatus
			err = workflow.ExecuteActivity(eksSelectFreeTrialDoNodesCtx, c.CreateSetupTopologyActivities.EksSelectFreeTrialNodes, params.Ou.OrgID).Get(eksSelectFreeTrialDoNodesCtx, &eksNodes)
			if err != nil {
				logger.Error("Failed to select eks free trial nodes", "Error", err)
				return err
			}
			for _, node := range eksNodes {
				var authCfg *authorized_clusters.K8sClusterConfig
				if node.ExtClusterCfgID > 0 {
					clusterAuthCtxKns := workflow.WithActivityOptions(ctx, ao)
					cerr := workflow.ExecuteActivity(clusterAuthCtxKns, c.CreateSetupTopologyActivities.GetClusterAuthCtxFromID, params.Ou, fmt.Sprintf("%d", node.ExtClusterCfgID)).Get(clusterAuthCtxKns, &authCfg)
					if cerr != nil {
						logger.Error("Failed to get cluster auth ctx", "Error", cerr)
						return cerr
					}
				}

				if authCfg != nil && !authCfg.IsPublic {
					if authCfg == nil || authCfg.Context == "" || authCfg.Region == "" {
						cerr := fmt.Errorf("failed to get cluster auth ctx")
						logger.Error("Failed to get cluster auth ctx", "Error", cerr)
						return cerr
					}
					eksDestroyNodePoolOrgResourcesCtx := workflow.WithActivityOptions(ctx, ao)
					logger.Info("Destroying node pool org resources", "PrivateEksRemoveNodePoolRequest", node)
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
			removeFreeTrialResourcesCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(removeFreeTrialResourcesCtx, c.CreateSetupTopologyActivities.RemoveFreeTrialOrgResources, params).Get(removeFreeTrialResourcesCtx, nil)
			if err != nil {
				logger.Error("Failed to add remove org free trial resources for account", "Error", err)
				return err
			}
			childWorkflowOptions := workflow.ChildWorkflowOptions{
				ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
				RetryPolicy:       ao.RetryPolicy,
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
				desWf := base_deploy_params.TopologyWorkflowRequest{
					TopologyDeployRequest: zeus_req_types.TopologyDeployRequest{
						TopologyID:                topID,
						CloudCtxNs:                params.CloudCtxNs,
						TopologyBaseInfraWorkload: *infraConfig,
					},
					OrgUser: params.Ou,
				}
				wfChildCtx := workflow.WithChildOptions(ctx, childWorkflowOptions)
				childWorkflowFuture := workflow.ExecuteChildWorkflow(wfChildCtx, "DestroyDeployedTopologyWorkflow", wfID, desWf)
				var childWE workflow.Execution
				if err = childWorkflowFuture.GetChildWorkflowExecution().Get(wfChildCtx, &childWE); err != nil {
					logger.Error("Failed to get child workflow execution", "Error", err)
					return err
				}
			}
		} else {
			updateResourcesToPaidCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(updateResourcesToPaidCtx, c.CreateSetupTopologyActivities.UpdateFreeTrialOrgResourcesToPaid, params).Get(updateResourcesToPaidCtx, nil)
			if err != nil {
				logger.Error("Failed to update org free trial resources to paid for account", "Error", err)
				return err
			}
		}
	}

	finishedCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
	if err != nil {
		logger.Error("Failed to update and mark orchestration inactive", "Error", err)
		return err
	}

	return nil
}
