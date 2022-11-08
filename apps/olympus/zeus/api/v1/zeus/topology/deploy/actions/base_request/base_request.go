package base_request

import (
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
)

type InternalDeploymentActionRequest struct {
	base_deploy_params.TopologyWorkflowRequest
}
