package deploy_workflow

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/cloud_ctx_logs"
	topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	deploy_topology_activities "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/create"
	deployment_status "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/status"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/temporal_actions/base_request"
	"go.temporal.io/sdk/temporal"
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

func (t *DeployTopologyWorkflow) GetWorkflows() []interface{} {
	return []interface{}{t.DeployTopologyWorkflow, t.DeployClusterTopologyWorkflow, t.DeployCronJobWorkflow, t.DeployJobWorkflow}
}

func (t *DeployTopologyWorkflow) DeployTopologyWorkflow(ctx workflow.Context, wfID string, params base_deploy_params.TopologyWorkflowRequest) error {
	logger := workflow.GetLogger(ctx)
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second * 60,
		BackoffCoefficient: 1.0,
		MaximumInterval:    time.Second * 60,
		MaximumAttempts:    100,
	}
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
		RetryPolicy:         retryPolicy,
	}
	oj := artemis_orchestrations.NewInternalActiveTemporalOrchestrationJobTemplate(wfID, "DeployTopologyWorkflow", "DeployTopologyWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	aerr := workflow.ExecuteActivity(alertCtx, "UpsertAssignmentV2", &oj).Get(alertCtx, &oj)
	if aerr != nil {
		logger.Error("Failed to upsert assignment", "Error", aerr)
		return aerr
	}

	t.DeployTopologyActivities.TopologyWorkflowRequest = params
	statusCtxKns := workflow.WithActivityOptions(ctx, ao)
	status := topology_deployment_status.NewPopulatedTopologyStatus(params.TopologyDeployRequest, topology_deployment_status.DeployInProgress)
	statusActivity := deployment_status.TopologyActivityDeploymentStatusActivity{}
	err := workflow.ExecuteActivity(statusCtxKns, statusActivity.CreateOrUpdateKubeCtxNsStatus, status.TopologyDeployRequest).Get(statusCtxKns, nil)
	if err != nil {
		logger.Error("Failed to update topology status", "Error", err)
		return err
	}

	statusCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(statusCtx, statusActivity.PostStatusUpdate, status.DeployStatus).Get(statusCtx, nil)
	if err != nil {
		logger.Error("Failed to update topology status", "Error", err)
		return err
	}

	deployParams := base_request.InternalDeploymentActionRequest{
		OrgUser: params.OrgUser,
		Kns:     params.TopologyDeployRequest,
	}
	ojl := &cloud_ctx_logs.CloudCtxNsLogs{
		OrchestrationID: oj.OrchestrationID,
		Status:          "Pending",
		Msg:             "DeployTopologyWorkflow starting",
		Ou:              params.OrgUser,
		CloudCtxNs:      deployParams.Kns.CloudCtxNs,
	}
	_ = StatusLogger(ctx, ao, ojl, "Pending", "DeployTopologyWorkflow starting")
	nsCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(nsCtx, t.DeployTopologyActivities.CreateNamespace, deployParams).Get(nsCtx, nil)
	if err != nil {
		_ = StatusLogger(ctx, ao, ojl, "Error", fmt.Sprintf("CreateNamespace failed: %e", err))
		logger.Error("Failed to create namespace", "Error", err)
		return err
	}

	chorCtx := workflow.WithActivityOptions(ctx, ao)
	if params.TopologyDeployRequest.RequestChoreographySecretDeploy == true {
		cmCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(cmCtx, t.DeployTopologyActivities.CreateChoreographySecret, deployParams).Get(chorCtx, nil)
		if err != nil {
			logger.Error("Failed to create choreographySecret", "Error", err)
			return err
		}
	}

	if params.TopologyDeployRequest.SecretRef != "" {
		secCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(secCtx, t.DeployTopologyActivities.CreateSecret, deployParams, params.TopologyDeployRequest.SecretRef).Get(secCtx, nil)
		if err != nil {
			logger.Error("Failed to get or deploy secret relationships", "Error", err)
			return err
		}
	}

	if params.TopologyDeployRequest.TopologyBaseInfraWorkload.ConfigMap != nil {
		_ = StatusLogger(ctx, ao, ojl, "Pending", "DeployConfigMap")
		cmCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(cmCtx, t.DeployTopologyActivities.DeployConfigMap, deployParams).Get(cmCtx, nil)
		if err != nil {
			_ = StatusLogger(ctx, ao, ojl, "Error", fmt.Sprintf("DeployConfigMap failed: %e", err))
			logger.Error("Failed to create configmap", "Error", err)
			return err
		} else {
			_ = StatusLogger(ctx, ao, ojl, "Success", "DeployConfigMap succeeded")
		}
	}

	if params.TopologyDeployRequest.TopologyBaseInfraWorkload.Deployment != nil {
		_ = StatusLogger(ctx, ao, ojl, "Pending", "DeployDeployment")
		logger.Error("Failed to create deployment", "Error", err)
		dCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(dCtx, t.DeployTopologyActivities.DeployDeployment, deployParams).Get(dCtx, nil)
		if err != nil {
			_ = StatusLogger(ctx, ao, ojl, "Error", fmt.Sprintf("DeployDeployment failed: %e", err))
			logger.Error("Failed to create deployment", "Error", err)
			return err
		} else {
			_ = StatusLogger(ctx, ao, ojl, "Success", "DeployDeployment succeeded")
		}
	}

	if params.TopologyDeployRequest.TopologyBaseInfraWorkload.StatefulSet != nil {
		_ = StatusLogger(ctx, ao, ojl, "Pending", "DeployStatefulSet")
		stsCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(stsCtx, t.DeployTopologyActivities.DeployStatefulSet, deployParams).Get(stsCtx, nil)
		if err != nil {
			_ = StatusLogger(ctx, ao, ojl, "Error", fmt.Sprintf("DeployStatefulSet failed: %e", err))
			logger.Error("Failed to create statefulset", "Error", err)
			return err
		} else {
			_ = StatusLogger(ctx, ao, ojl, "Success", "DeployStatefulSet succeeded")
		}
	}

	if params.TopologyDeployRequest.TopologyBaseInfraWorkload.Service != nil {
		_ = StatusLogger(ctx, ao, ojl, "Pending", "DeployService")
		svcCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(svcCtx, t.DeployTopologyActivities.DeployService, deployParams).Get(svcCtx, nil)
		if err != nil {
			_ = StatusLogger(ctx, ao, ojl, "Error", fmt.Sprintf("DeployService failed: %e", err))
			logger.Error("Failed to create service", "Error", err)
			return err
		} else {
			_ = StatusLogger(ctx, ao, ojl, "Success", "DeployService succeeded")
		}
	}

	if params.TopologyDeployRequest.TopologyBaseInfraWorkload.Ingress != nil {
		_ = StatusLogger(ctx, ao, ojl, "Pending", "DeployIngress")
		ingCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(ingCtx, t.DeployTopologyActivities.DeployIngress, deployParams).Get(ingCtx, nil)
		if err != nil {
			_ = StatusLogger(ctx, ao, ojl, "Error", fmt.Sprintf("DeployIngress failed: %e", err))
			logger.Error("Failed to create ingress", "Error", err)
			return err
		} else {
			_ = StatusLogger(ctx, ao, ojl, "Success", "DeployIngress succeeded")
		}
	}
	if params.TopologyDeployRequest.TopologyBaseInfraWorkload.ServiceMonitor != nil {
		_ = StatusLogger(ctx, ao, ojl, "Pending", "DeployServiceMonitor")
		smCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(smCtx, t.DeployTopologyActivities.DeployServiceMonitor, deployParams).Get(smCtx, nil)
		if err != nil {
			_ = StatusLogger(ctx, ao, ojl, "Error", fmt.Sprintf("DeployServiceMonitor failed: %e", err))
			logger.Error("Failed to create servicemonitor", "Error", err)
			return err
		} else {
			_ = StatusLogger(ctx, ao, ojl, "Success", "DeployServiceMonitor succeeded")
		}
	}
	status.TopologyStatus = topology_deployment_status.DeployComplete
	err = workflow.ExecuteActivity(statusCtx, statusActivity.PostStatusUpdate, status.DeployStatus).Get(statusCtx, nil)
	if err != nil {
		logger.Error("Failed to update topology status", "Error", err)
		return err
	}
	_ = StatusLogger(ctx, ao, ojl, "Success", "DeployTopologyWorkflow succeeded")
	finishedCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
	if err != nil {
		logger.Error("Failed to update and mark orchestration inactive", "Error", err)
		return err
	}
	return err
}

func StatusLogger(ctx workflow.Context, ao workflow.ActivityOptions, ojl *cloud_ctx_logs.CloudCtxNsLogs, status, msg string) error {
	ojl.Status = status
	ojl.Msg = msg
	logCtx := workflow.WithActivityOptions(ctx, ao)
	lerr := workflow.ExecuteActivity(logCtx, "InsertClusterLogs", ojl).Get(logCtx, nil)
	if lerr != nil {
		log.Err(lerr).Interface("oj", ojl).Msg("UpsertAssignment: InsertCloudCtxNsLog failed")
		return lerr
	}
	return nil
}
