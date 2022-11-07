package deploy_topology

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	deploy_topology "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/create"
	"go.temporal.io/sdk/workflow"
)

//func WorkflowDeploy(ctx workflow.Context, params WorkflowParams) error {
//	var activities *YourActivityStruct
//	future := workflow.ExecuteActivity(ctx, activities.Activity, ctx)
//	var yourActivityResult YourActivityResult
//	if err := future.Get(ctx, &yourActivityResult); err != nil {
//		// ...
//	}
//	return nil
//}

type DeployTopologyParams struct {
	Kns        zeus_core.KubeCtxNs
	TopologyID int
	UserID     int
	OrgID      int
	chart_workload.NativeK8s
}

type DeployTopologyWorkflow struct {
	temporal_base.Workflow
	deploy_topology.DeployTopologyActivity
}

const defaultTimeout = 3 * time.Minute

func (t *DeployTopologyWorkflow) DeployTopologyWorkflow(ctx workflow.Context, params DeployTopologyParams) error {
	log := workflow.GetLogger(ctx)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}

	nsCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(nsCtx, t.DeployTopologyActivity.CreateNamespace, params).Get(nsCtx, nil)
	if err != nil {
		log.Error("Failed to create namespace", "Error", err)
		return err
	}

	if params.Deployment != nil {
		dCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(dCtx, t.DeployTopologyActivity.DeployDeployment, params).Get(dCtx, nil)
		if err != nil {
			log.Error("Failed to create deployment", "Error", err)
			return err
		}
	}

	if params.StatefulSet != nil {
		stsCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(stsCtx, t.DeployTopologyActivity.DeployStatefulSet, params).Get(stsCtx, nil)
		if err != nil {
			log.Error("Failed to create statefulset", "Error", err)
			return err
		}
	}

	if params.Service != nil {
		svcCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(svcCtx, t.DeployTopologyActivity.DeployService, params).Get(svcCtx, nil)
		if err != nil {
			log.Error("Failed to create service", "Error", err)
			return err
		}
	}

	if params.Ingress != nil {
		ingCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(ingCtx, t.DeployTopologyActivity.DeployIngress, params).Get(ingCtx, nil)
		if err != nil {
			log.Error("Failed to create ingress", "Error", err)
			return err
		}
	}
	return err
}
