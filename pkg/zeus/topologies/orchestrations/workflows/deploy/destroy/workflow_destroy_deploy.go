package destroy_deployed_workflow

import (
	"time"

	topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	destroy_deploy_activities "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/destroy"
	deployment_status "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/status"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/actions/base_request"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/actions/deploy/workload_state"
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

	deployParams := base_request.InternalDeploymentActionRequest{
		Kns:       params.Kns,
		OrgUser:   params.OrgUser,
		NativeK8s: params.NativeK8s,
	}
	if params.Deployment != nil {
		dCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(dCtx, t.DestroyDeployTopologyActivities.DestroyDeployDeployment, deployParams).Get(dCtx, nil)
		if err != nil {
			log.Error("Failed to destroy deployment", "Error", err)
			return err
		}
	}

	if params.StatefulSet != nil {
		stsCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(stsCtx, t.DestroyDeployTopologyActivities.DestroyDeployStatefulSet, deployParams).Get(stsCtx, nil)
		if err != nil {
			log.Error("Failed to destroy statefulset", "Error", err)
			return err
		}
	}

	if params.Service != nil {
		svcCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(svcCtx, t.DestroyDeployTopologyActivities.DestroyDeployService, deployParams).Get(svcCtx, nil)
		if err != nil {
			log.Error("Failed to destroy service", "Error", err)
			return err
		}
	}

	if params.Ingress != nil {
		ingCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(ingCtx, t.DestroyDeployTopologyActivities.DestroyDeployIngress, deployParams).Get(ingCtx, nil)
		if err != nil {
			log.Error("Failed to destroy ingress", "Error", err)
			return err
		}
	}
	statusActivity.TopologyStatus = topology_deployment_status.Complete
	err = workflow.ExecuteActivity(statusCtx, statusActivity.PostStatusUpdate, statusActivity.InternalWorkloadStatusUpdateRequest).Get(statusCtx, nil)
	if err != nil {
		log.Error("Failed to update topology status", "Error", err)
		return err
	}
	return nil
}
