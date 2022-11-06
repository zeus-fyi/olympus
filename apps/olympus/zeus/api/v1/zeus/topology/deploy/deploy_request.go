package deploy

import (
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/external/create_or_update"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/external/delete"
)

type DeploymentActionRequest struct {
	Action string
	create_or_update_deploy.TopologyDeployCreateActionDeployRequest
	delete_deploy.TopologyDeployActionDeleteDeploymentRequest
}
