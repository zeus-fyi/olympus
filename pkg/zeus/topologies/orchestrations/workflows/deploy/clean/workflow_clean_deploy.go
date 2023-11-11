package clean_deployed_workflow

import (
	"time"

	topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	destroy_deploy_activities "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/destroy"
	deployment_status "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/status"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/temporal_actions/base_request"
	"go.temporal.io/sdk/workflow"
)

type CleanDeployTopologyWorkflow struct {
	temporal_base.Workflow
	destroy_deploy_activities.DestroyDeployTopologyActivities
}

const defaultTimeout = 10 * time.Minute

func (t *CleanDeployTopologyWorkflow) GetWorkflow() interface{} {
	return t.CleanDeployedTopologyWorkflow
}

func NewCleanDeployTopologyWorkflow() CleanDeployTopologyWorkflow {
	cleanDeployWf := CleanDeployTopologyWorkflow{
		Workflow:                        temporal_base.Workflow{},
		DestroyDeployTopologyActivities: destroy_deploy_activities.DestroyDeployTopologyActivities{},
	}
	return cleanDeployWf
}

func (t *CleanDeployTopologyWorkflow) CleanDeployedTopologyWorkflow(ctx workflow.Context, params base_deploy_params.TopologyWorkflowRequest) error {
	log := workflow.GetLogger(ctx)

	t.DestroyDeployTopologyActivities.TopologyWorkflowRequest = params
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}

	statusCtx := workflow.WithActivityOptions(ctx, ao)
	status := topology_deployment_status.NewPopulatedTopologyStatus(params.TopologyDeployRequest, topology_deployment_status.CleanDeployInProgress)
	statusActivity := deployment_status.TopologyActivityDeploymentStatusActivity{}

	err := workflow.ExecuteActivity(statusCtx, statusActivity.PostStatusUpdate, status.DeployStatus).Get(statusCtx, nil)
	if err != nil {
		log.Error("Failed to update topology status", "Error", err)
		return err
	}

	deployParams := base_request.InternalDeploymentActionRequest{
		Kns:     params.TopologyDeployRequest,
		OrgUser: params.OrgUser,
	}

	if params.TopologyDeployRequest.TopologyBaseInfraWorkload.ConfigMap != nil {
		cmCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(cmCtx, t.DestroyDeployTopologyActivities.DestroyDeployConfigMap, deployParams).Get(cmCtx, nil)
		if err != nil {
			log.Error("Failed to destroy configmap", "Error", err)
			return err
		}
	}

	if params.TopologyDeployRequest.TopologyBaseInfraWorkload.Deployment != nil {
		dCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(dCtx, t.DestroyDeployTopologyActivities.DestroyDeployDeployment, deployParams).Get(dCtx, nil)
		if err != nil {
			log.Error("Failed to destroy deployment", "Error", err)
			return err
		}
	}

	if params.TopologyDeployRequest.TopologyBaseInfraWorkload.StatefulSet != nil {
		stsCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(stsCtx, t.DestroyDeployTopologyActivities.DestroyDeployStatefulSet, deployParams).Get(stsCtx, nil)
		if err != nil {
			log.Error("Failed to destroy statefulset", "Error", err)
			return err
		}
	}

	if params.TopologyDeployRequest.TopologyBaseInfraWorkload.Service != nil {
		svcCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(svcCtx, t.DestroyDeployTopologyActivities.DestroyDeployService, deployParams).Get(svcCtx, nil)
		if err != nil {
			log.Error("Failed to destroy service", "Error", err)
			return err
		}
	}

	if params.TopologyDeployRequest.TopologyBaseInfraWorkload.Ingress != nil {
		ingCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(ingCtx, t.DestroyDeployTopologyActivities.DestroyDeployIngress, deployParams).Get(ingCtx, nil)
		if err != nil {
			log.Error("Failed to destroy ingress", "Error", err)
			return err
		}
	}

	knsCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(knsCtx, statusActivity.DeleteKubeCtxNsStatus, status.TopologyDeployRequest).Get(knsCtx, nil)
	if err != nil {
		log.Error("Failed to remove topology kns status", "Error", err)
		return err
	}

	status.TopologyStatus = topology_deployment_status.CleanDeployComplete
	err = workflow.ExecuteActivity(statusCtx, statusActivity.PostStatusUpdate, status.DeployStatus).Get(statusCtx, nil)
	if err != nil {
		log.Error("Failed to update topology status", "Error", err)
		return err
	}
	return nil
}
