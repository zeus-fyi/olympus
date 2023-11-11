package deploy_workflow

import (
	"time"

	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/temporal_actions/base_request"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (t *DeployTopologyWorkflow) DeployCronJobWorkflow(ctx workflow.Context, params base_deploy_params.TopologyWorkflowRequest) error {
	log := workflow.GetLogger(ctx)

	t.DeployTopologyActivities.TopologyWorkflowRequest = params
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second * 60,
		BackoffCoefficient: 1.0,
		MaximumInterval:    time.Second * 60,
	}
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
		RetryPolicy:         retryPolicy,
	}
	deployParams := base_request.InternalDeploymentActionRequest{
		Kns:                       params.TopologyDeployRequest,
		OrgUser:                   params.OrgUser,
		TopologyBaseInfraWorkload: params.TopologyBaseInfraWorkload,
		ClusterName:               params.ClusterClassName,
		SecretRef:                 params.SecretRef,
	}
	nsCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(nsCtx, t.DeployTopologyActivities.CreateNamespace, deployParams).Get(nsCtx, nil)
	if err != nil {
		log.Error("Failed to create namespace", "Error", err)
		return err
	}

	if params.ConfigMap != nil {
		cmCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(cmCtx, t.DeployTopologyActivities.DeployConfigMap, deployParams).Get(cmCtx, nil)
		if err != nil {
			log.Error("Failed to create configmap", "Error", err)
			return err
		}
	}

	if params.SecretRef != "" {
		secCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(secCtx, t.DeployTopologyActivities.CreateSecret, deployParams).Get(secCtx, nil)
		if err != nil {
			log.Error("Failed to get or deploy secret relationships", "Error", err)
			return err
		}
	}

	if params.CronJob != nil {
		cjCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(cjCtx, t.DeployTopologyActivities.CreateCronJob, deployParams).Get(cjCtx, nil)
		if err != nil {
			log.Error("Failed to create cronjob", "Error", err)
			return err
		}
	}

	if params.Service != nil {
		svcCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(svcCtx, t.DeployTopologyActivities.DeployService, deployParams).Get(svcCtx, nil)
		if err != nil {
			log.Error("Failed to create service", "Error", err)
			return err
		}
	}
	return nil
}
