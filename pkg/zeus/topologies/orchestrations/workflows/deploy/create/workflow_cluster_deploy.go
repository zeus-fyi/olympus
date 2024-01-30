package deploy_workflow

import (
	"errors"
	"fmt"
	"time"

	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"github.com/zeus-fyi/zeus/zeus/workload_config_drivers/topology_workloads"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (t *DeployTopologyWorkflow) DeployClusterTopologyWorkflow(ctx workflow.Context, wfID string, params base_deploy_params.ClusterTopologyWorkflowRequest) error {
	logger := workflow.GetLogger(ctx)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 60,
			BackoffCoefficient: 1.0,
			MaximumInterval:    time.Second * 60,
			MaximumAttempts:    1000,
		},
	}

	for _, topID := range params.TopologyIDs {
		req := zeus_req_types.TopologyDeployRequest{
			TopologyID:                      topID,
			ClusterClassName:                params.ClusterClassName,
			CloudCtxNs:                      params.CloudCtxNs,
			RequestChoreographySecretDeploy: params.RequestChoreographySecret,
		}
		if req.Context == "" || req.Namespace == "" || req.Region == "" || req.CloudProvider == "" {
			return errors.New("cloudCtxNs is empty")
		}
		if params.AppTaint {
			req.ClusterClassName = params.ClusterClassName
		}
		if params.ClusterClassName != "" {
			req.SecretRef = params.ClusterClassName
		}

		var infraConfig *topology_workloads.TopologyBaseInfraWorkload
		deployStatusCtx := workflow.WithActivityOptions(ctx, ao)
		err := workflow.ExecuteActivity(deployStatusCtx, t.DeployTopologyActivities.GetTopologyInfraConfig, params.OrgUser, topID).Get(deployStatusCtx, &infraConfig)
		if err != nil {
			logger.Error("Failed to get topology infra config", "Error", err)
			return err
		}
		if infraConfig == nil {
			err = fmt.Errorf("infraConfig is nil")
			logger.Error("Failed to get topology infra config", "Error", err)
			return err
		}
		topParams := zeus_req_types.TopologyDeployRequest{
			TopologyID:                      topID,
			ClusterClassName:                params.ClusterClassName,
			CloudCtxNs:                      params.CloudCtxNs,
			SecretRef:                       req.SecretRef,
			RequestChoreographySecretDeploy: params.RequestChoreographySecret,
			TopologyBaseInfraWorkload:       *infraConfig, // nil check is above
		}
		deployChildWorkflowOptions := workflow.ChildWorkflowOptions{
			ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
			RetryPolicy:       ao.RetryPolicy,
		}
		topWfReq := base_deploy_params.TopologyWorkflowRequest{
			TopologyDeployRequest: topParams,
			OrgUser:               params.OrgUser,
		}
		clusterDeployCtx := workflow.WithChildOptions(ctx, deployChildWorkflowOptions)
		deployChildWorkflowFuture := workflow.ExecuteChildWorkflow(clusterDeployCtx, "DeployTopologyWorkflow", wfID, topWfReq)
		var deployChildWfExec workflow.Execution
		if err = deployChildWorkflowFuture.GetChildWorkflowExecution().Get(ctx, &deployChildWfExec); err != nil {
			logger.Error("Failed to get child deployment workflow execution", "Error", err)
			return err
		}
	}
	return nil
}
