package destroy_deployed_workflow

import (
	"time"

	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/temporal_actions/base_request"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (t *DestroyDeployTopologyWorkflow) DestroyCronJobWorkflow(ctx workflow.Context, params base_deploy_params.TopologyWorkflowRequest) error {
	log := workflow.GetLogger(ctx)

	t.DestroyDeployTopologyActivities.TopologyWorkflowRequest = params
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second * 60,
		BackoffCoefficient: 1.0,
		MaximumInterval:    time.Minute * 60,
	}
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
		RetryPolicy:         retryPolicy,
	}

	deployParams := base_request.InternalDeploymentActionRequest{
		Kns:                       params.TopologyDeployRequest,
		OrgUser:                   params.OrgUser,
		TopologyBaseInfraWorkload: params.TopologyBaseInfraWorkload,
	}

	if params.ConfigMap != nil {
		cmCtx := workflow.WithActivityOptions(ctx, ao)
		err := workflow.ExecuteActivity(cmCtx, t.DestroyDeployTopologyActivities.DestroyDeployConfigMap, deployParams).Get(cmCtx, nil)
		if err != nil {
			log.Error("Failed to destroy configmap", "Error", err)
			return err
		}
	}

	if params.Service != nil {
		svcCtx := workflow.WithActivityOptions(ctx, ao)
		err := workflow.ExecuteActivity(svcCtx, t.DestroyDeployTopologyActivities.DestroyDeployService, deployParams).Get(svcCtx, nil)
		if err != nil {
			log.Error("Failed to destroy service", "Error", err)
			return err
		}
	}

	if params.CronJob != nil {
		ingCtx := workflow.WithActivityOptions(ctx, ao)
		err := workflow.ExecuteActivity(ingCtx, t.DestroyDeployTopologyActivities.DestroyCronJob, deployParams).Get(ingCtx, nil)
		if err != nil {
			log.Error("Failed to destroy job", "Error", err)
			return err
		}
	}

	//if params.ServiceMonitor != nil {
	//	smCtx := workflow.WithActivityOptions(ctx, ao)
	//	err := workflow.ExecuteActivity(smCtx, t.DestroyDeployTopologyActivities.DestroyDeployServiceMonitor, deployParams).Get(smCtx, nil)
	//	if err != nil {
	//		log.Error("Failed to destroy servicemonitor", "Error", err)
	//		return err
	//	}
	//}

	return nil
}
