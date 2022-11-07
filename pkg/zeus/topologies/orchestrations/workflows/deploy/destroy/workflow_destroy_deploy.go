package destroy_deployed_workflow

import (
	"time"

	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	destroy_deploy "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/destroy"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"go.temporal.io/sdk/workflow"
)

type DestroyDeployTopologyWorkflow struct {
	temporal_base.Workflow
	destroy_deploy.DestroyDeployTopologyActivity
}

const defaultTimeout = 3 * time.Minute

func (t *DestroyDeployTopologyWorkflow) DestroyDeployedTopologyWorkflow(ctx workflow.Context, params base_deploy_params.DeployTopologyParams) error {
	log := workflow.GetLogger(ctx)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}

	if params.Deployment != nil {
		dCtx := workflow.WithActivityOptions(ctx, ao)
		err := workflow.ExecuteActivity(dCtx, t.DestroyDeployTopologyActivity.DestroyDeployDeployment, params).Get(dCtx, nil)
		if err != nil {
			log.Error("Failed to destroy deployment", "Error", err)
			return err
		}
	}

	if params.StatefulSet != nil {
		stsCtx := workflow.WithActivityOptions(ctx, ao)
		err := workflow.ExecuteActivity(stsCtx, t.DestroyDeployTopologyActivity.DestroyDeployStatefulSet, params).Get(stsCtx, nil)
		if err != nil {
			log.Error("Failed to destroy statefulset", "Error", err)
			return err
		}
	}

	if params.Service != nil {
		svcCtx := workflow.WithActivityOptions(ctx, ao)
		err := workflow.ExecuteActivity(svcCtx, t.DestroyDeployTopologyActivity.DestroyDeployService, params).Get(svcCtx, nil)
		if err != nil {
			log.Error("Failed to destroy service", "Error", err)
			return err
		}
	}

	if params.Ingress != nil {
		ingCtx := workflow.WithActivityOptions(ctx, ao)
		err := workflow.ExecuteActivity(ingCtx, t.DestroyDeployTopologyActivity.DestroyDeployIngress, params).Get(ingCtx, nil)
		if err != nil {
			log.Error("Failed to destroy ingress", "Error", err)
			return err
		}
	}
	return nil
}
