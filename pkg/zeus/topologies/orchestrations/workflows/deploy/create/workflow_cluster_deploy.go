package deploy_workflow

import (
	"errors"

	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"go.temporal.io/sdk/workflow"
)

func (t *DeployTopologyWorkflow) DeployClusterTopologyWorkflow(ctx workflow.Context, wfID string, params base_deploy_params.ClusterTopologyWorkflowRequest) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	deployStatusCtx := workflow.WithActivityOptions(ctx, ao)
	for _, topID := range params.TopologyIDs {
		req := zeus_req_types.TopologyDeployRequest{
			TopologyID:                      topID,
			CloudCtxNs:                      params.CloudCtxNS,
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
		err := workflow.ExecuteActivity(deployStatusCtx, t.DeployTopologyActivities.DeployClusterTopology, req, params.OrgUser).Get(deployStatusCtx, nil)
		if err != nil {
			log.Error("Failed to deploy topology from cluster definition", "Error", err)
			return err
		}
	}
	return nil
}
