package deploy_workflow

import (
	"time"

	topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	deploy_topology_activities "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/create"
	deployment_status "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/status"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/actions/deploy/workload_state"
	"go.temporal.io/sdk/workflow"
)

type DeployTopologyWorkflow struct {
	temporal_base.Workflow
	deploy_topology_activities.DeployTopologyActivities
}

const defaultTimeout = 10 * time.Minute

func NewDeployTopologyWorkflow() DeployTopologyWorkflow {
	deployWf := DeployTopologyWorkflow{
		Workflow:                 temporal_base.Workflow{},
		DeployTopologyActivities: deploy_topology_activities.DeployTopologyActivities{},
	}
	return deployWf
}

func (t *DeployTopologyWorkflow) GetWorkflow() interface{} {
	return t.DeployTopologyWorkflow
}

func (t *DeployTopologyWorkflow) DeployTopologyWorkflow(ctx workflow.Context, params base_deploy_params.TopologyWorkflowRequest) error {
	log := workflow.GetLogger(ctx)

	t.DeployTopologyActivities.TopologyWorkflowRequest = params
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}

	statusCtx := workflow.WithActivityOptions(ctx, ao)
	statusActivity := deployment_status.TopologyActivityDeploymentStatusActivity{
		Host: params.Host,
		InternalWorkloadStatusUpdateRequest: workload_state.InternalWorkloadStatusUpdateRequest{
			TopologyID:     params.Kns.TopologyID,
			TopologyStatus: topology_deployment_status.InProgress,
		},
	}
	err := workflow.ExecuteActivity(statusCtx, statusActivity.PostStatusUpdate, statusActivity.InternalWorkloadStatusUpdateRequest).Get(statusCtx, nil)
	if err != nil {
		log.Error("Failed to update topology status", "Error", err)
		return err
	}

	nsCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(nsCtx, t.DeployTopologyActivities.CreateNamespace).Get(nsCtx, nil)
	if err != nil {
		log.Error("Failed to create namespace", "Error", err)
		return err
	}

	if params.Deployment != nil {
		dCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(dCtx, t.DeployTopologyActivities.DeployDeployment).Get(dCtx, nil)
		if err != nil {
			log.Error("Failed to create deployment", "Error", err)
			return err
		}
	}

	if params.StatefulSet != nil {
		stsCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(stsCtx, t.DeployTopologyActivities.DeployStatefulSet).Get(stsCtx, nil)
		if err != nil {
			log.Error("Failed to create statefulset", "Error", err)
			return err
		}
	}

	if params.Service != nil {
		svcCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(svcCtx, t.DeployTopologyActivities.DeployService).Get(svcCtx, nil)
		if err != nil {
			log.Error("Failed to create service", "Error", err)
			return err
		}
	}

	if params.Ingress != nil {
		ingCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(ingCtx, t.DeployTopologyActivities.DeployIngress, params).Get(ingCtx, nil)
		if err != nil {
			log.Error("Failed to create ingress", "Error", err)
			return err
		}
	}

	statusActivity.TopologyStatus = topology_deployment_status.Complete
	err = workflow.ExecuteActivity(statusCtx, statusActivity.PostStatusUpdate, statusActivity.InternalWorkloadStatusUpdateRequest).Get(statusCtx, nil)
	if err != nil {
		log.Error("Failed to update topology status", "Error", err)
		return err
	}

	return err
}
