package deploy

import (
	create_or_update_deploy "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/create_or_update"
	delete_deploy "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/delete"
)

type DeploymentActionRequest struct {
	Action string
	create_or_update_deploy.TopologyDeployCreateActionDeployRequest
	delete_deploy.TopologyDeployActionDeleteDeploymentRequest
}
