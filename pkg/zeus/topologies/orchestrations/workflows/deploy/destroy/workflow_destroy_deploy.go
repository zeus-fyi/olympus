package destroy_deployed_workflow

import (
	"time"

	topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
	create_topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/definitions/state"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	destroy_deploy_activities "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/destroy"
	deployment_status "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/status"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"go.temporal.io/sdk/workflow"
)

type DestroyDeployTopologyWorkflow struct {
	temporal_base.Workflow
	destroy_deploy_activities.DestroyDeployTopologyActivities
}

const defaultTimeout = 10 * time.Minute

func (t *DestroyDeployTopologyWorkflow) GetWorkflow() interface{} {
	return t.DestroyDeployedTopologyWorkflow
}

func NewDestroyDeployTopologyWorkflow() DestroyDeployTopologyWorkflow {
	destroyDeployWf := DestroyDeployTopologyWorkflow{
		Workflow:                        temporal_base.Workflow{},
		DestroyDeployTopologyActivities: destroy_deploy_activities.DestroyDeployTopologyActivities{},
	}
	return destroyDeployWf
}

func (t *DestroyDeployTopologyWorkflow) DestroyDeployedTopologyWorkflow(ctx workflow.Context, params base_deploy_params.TopologyWorkflowRequest) error {
	log := workflow.GetLogger(ctx)

	t.DestroyDeployTopologyActivities.TopologyWorkflowRequest = params
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}

	statusCtx := workflow.WithActivityOptions(ctx, ao)
	statusActivity := deployment_status.TopologyActivityDeploymentStatusActivity{
		Host:             params.Host,
		DeploymentStatus: create_topology_deployment_status.DeploymentStatus{},
	}
	statusActivity.Status.TopologiesDeployed.TopologyID = params.Kns.TopologyID
	statusActivity.Status.TopologiesDeployed.TopologyStatus = topology_deployment_status.InProgress
	err := workflow.ExecuteActivity(statusCtx, statusActivity.PostStatusUpdate).Get(statusCtx, nil)
	if err != nil {
		log.Error("Failed to update topology status", "Error", err)
		return err
	}

	if params.Deployment != nil {
		dCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(dCtx, t.DestroyDeployTopologyActivities.DestroyDeployDeployment).Get(dCtx, nil)
		if err != nil {
			log.Error("Failed to destroy deployment", "Error", err)
			return err
		}
	}

	if params.StatefulSet != nil {
		stsCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(stsCtx, t.DestroyDeployTopologyActivities.DestroyDeployStatefulSet).Get(stsCtx, nil)
		if err != nil {
			log.Error("Failed to destroy statefulset", "Error", err)
			return err
		}
	}

	if params.Service != nil {
		svcCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(svcCtx, t.DestroyDeployTopologyActivities.DestroyDeployService).Get(svcCtx, nil)
		if err != nil {
			log.Error("Failed to destroy service", "Error", err)
			return err
		}
	}

	if params.Ingress != nil {
		ingCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(ingCtx, t.DestroyDeployTopologyActivities.DestroyDeployIngress).Get(ingCtx, nil)
		if err != nil {
			log.Error("Failed to destroy ingress", "Error", err)
			return err
		}
	}
	statusActivity.Status.TopologiesDeployed.TopologyStatus = topology_deployment_status.Complete
	err = workflow.ExecuteActivity(statusCtx, statusActivity.PostStatusUpdate).Get(statusCtx, nil)
	if err != nil {
		log.Error("Failed to update topology status", "Error", err)
		return err
	}
	return nil
}
