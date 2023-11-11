package deploy_workflow

import (
	"errors"
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
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

		var infraConfig *chart_workload.TopologyBaseInfraWorkload
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
		topParams := base_deploy_params.TopologyWorkflowRequest{
			Kns: kns.TopologyKubeCtxNs{
				TopologyID: topID,
				CloudCtxNs: params.CloudCtxNs,
			},
			OrgUser:                   params.OrgUser,
			RequestChoreographySecret: params.RequestChoreographySecret,
			ClusterClassName:          params.ClusterClassName,
			SecretRef:                 req.SecretRef,
			TopologyBaseInfraWorkload: *infraConfig, // nil check is above
		}
		deployChildWorkflowOptions := workflow.ChildWorkflowOptions{
			ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
		}
		clusterDeployCtx := workflow.WithChildOptions(ctx, deployChildWorkflowOptions)
		deployChildWorkflowFuture := workflow.ExecuteChildWorkflow(clusterDeployCtx, "DeployTopologyWorkflow", wfID, topParams)
		var deployChildWfExec workflow.Execution
		if err = deployChildWorkflowFuture.GetChildWorkflowExecution().Get(ctx, &deployChildWfExec); err != nil {
			logger.Error("Failed to get child deployment workflow execution", "Error", err)
			return err
		}
	}
	return nil
}
