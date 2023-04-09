package deploy_workflow_cluster_setup

import (
	"time"

	"github.com/google/uuid"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	deploy_topology_activities_create_setup "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/cluster_setup"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
	"go.temporal.io/sdk/workflow"
)

type ClusterSetupWorkflow struct {
	temporal_base.Workflow
	deploy_topology_activities_create_setup.CreateSetupTopologyActivities
}

const defaultTimeout = 10 * time.Minute

func NewDeployCreateSetupTopologyWorkflow() ClusterSetupWorkflow {
	deployWf := ClusterSetupWorkflow{
		Workflow:                      temporal_base.Workflow{},
		CreateSetupTopologyActivities: deploy_topology_activities_create_setup.CreateSetupTopologyActivities{},
	}
	return deployWf
}

func (c *ClusterSetupWorkflow) GetWorkflows() []interface{} {
	return []interface{}{c.DeployClusterSetupWorkflow}
}

type ClusterSetupRequest struct {
	Ou org_users.OrgUser
	zeus_common_types.CloudCtxNs
	ClusterID     uuid.UUID
	Nodes         hestia_autogen_bases.Nodes
	NodesQuantity float64
	Disks         hestia_autogen_bases.Disks
	DisksQuantity float64
}

func (c *ClusterSetupWorkflow) DeployClusterSetupWorkflow(ctx workflow.Context, params ClusterSetupRequest) error {
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
	err = workflow.ExecuteActivity(diskOrgResourcesCtx, c.CreateSetupTopologyActivities.AddDiskResourcesToOrg, params).Get(diskOrgResourcesCtx, nil)
	if err != nil {
		log.Error("Failed to add disk resources to org account", "Error", err)
		return err
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
	return err
}
