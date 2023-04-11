package deploy_workflow_cluster_setup

import (
	"context"
	"time"

	hestia_stripe "github.com/zeus-fyi/olympus/pkg/hestia/stripe"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	deploy_topology_activities_create_setup "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/cluster_setup"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
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

type NodePoolRequestStatus struct {
	ClusterID  string
	NodePoolID string
}

func (c *ClusterSetupWorkflow) DeployClusterSetupWorkflow(ctx workflow.Context, params base_deploy_params.ClusterSetupRequest) error {
	log := workflow.GetLogger(ctx)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}

	// TODO add billing email step
	nodePoolRequestStatusCtxKns := workflow.WithActivityOptions(ctx, ao)
	var nodePoolRequestStatus NodePoolRequestStatus
	err := workflow.ExecuteActivity(nodePoolRequestStatusCtxKns, c.CreateSetupTopologyActivities.MakeNodePoolRequest, params).Get(nodePoolRequestStatusCtxKns, &nodePoolRequestStatus)
	if err != nil {
		log.Error("Failed to complete node pool request", "Error", err)
		return err
	}

	nodePoolOrgResourcesCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(nodePoolRequestStatusCtxKns, c.CreateSetupTopologyActivities.AddNodePoolToOrgResources, params, nodePoolRequestStatus).Get(nodePoolOrgResourcesCtx, nil)
	if err != nil {
		log.Error("Failed to add node resources to org account", "Error", err)
		return err
	}

	authCloudCtxNsCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(authCloudCtxNsCtx, c.CreateSetupTopologyActivities.AddAuthCtxNsOrg, params).Get(authCloudCtxNsCtx, nil)
	if err != nil {
		log.Error("Failed to authorize auth ns to org account", "Error", err)
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
	params.CloudCtxNs.Namespace = params.ClusterID.String()
	clusterDeployCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(clusterDeployCtx, c.CreateSetupTopologyActivities.DeployClusterTopologyFromUI, params.Cluster.ClusterName, sbNames, params.CloudCtxNs, params.Ou).Get(clusterDeployCtx, nil)
	if err != nil {
		log.Error("Failed to add deploy cluster", "Error", err)
		return err
	}

	if params.FreeTrial {
		err = workflow.Sleep(ctx, 60*time.Hour)
		if err != nil {
			log.Error("Failed to sleep for 1 hour", "Error", err)
			return err
		}
		hestiaCtx := context.Background()
		isBillingSetup, herr := hestia_stripe.DoesUserHaveBillingMethod(hestiaCtx, params.Ou.UserID)
		if herr != nil {
			log.Error("Failed to check if user has billing method", "Error", herr)
			return herr
		}
		if !isBillingSetup {
			destroyClusterCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(destroyClusterCtx, c.CreateSetupTopologyActivities.DestroyCluster, params.CloudCtxNs).Get(destroyClusterCtx, nil)
			if err != nil {
				log.Error("Failed to add deploy cluster", "Error", err)
				return err
			}
			removeAuthCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(removeAuthCtx, c.CreateSetupTopologyActivities.RemoveAuthCtxNsOrg, params).Get(removeAuthCtx, nil)
			if err != nil {
				log.Error("Failed to add deploy cluster", "Error", err)
				return err
			}
			destroyNodePoolOrgResourcesCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(destroyNodePoolOrgResourcesCtx, c.CreateSetupTopologyActivities.RemoveNodePoolRequest, params, destroyNodePoolOrgResourcesCtx).Get(destroyNodePoolOrgResourcesCtx, nil)
			if err != nil {
				log.Error("Failed to add remove node resources for account", "Error", err)
				return err
			}
			removeFreeTrialResourcesCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(removeFreeTrialResourcesCtx, c.CreateSetupTopologyActivities.RemoveFreeTrialOrgResources, params).Get(removeFreeTrialResourcesCtx, nil)
			if err != nil {
				log.Error("Failed to add remove org free trial resources for account", "Error", err)
				return err
			}
		} else {
			updateResourcesToPaidCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(updateResourcesToPaidCtx, c.CreateSetupTopologyActivities.UpdateFreeTrialOrgResourcesToPaid, params).Get(updateResourcesToPaidCtx, nil)
			if err != nil {
				log.Error("Failed to update org free trial resources to paid for account", "Error", err)
				return err
			}
		}
	}
	return err
}
