package base_request

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

type InternalDeploymentActionRequest struct {
	Kns     zeus_req_types.TopologyDeployRequest `json:"topologyDeployRequest"`
	OrgUser org_users.OrgUser                    `json:"orgUser"`
}

type ClusterDeployActionRequest struct {
	Kns     zeus_req_types.TopologyDeployRequest `json:"topologyDeployRequest"`
	OrgUser org_users.OrgUser                    `json:"orgUser"`
}

type ExternalDeploymentActionRequest struct {
	zeus_req_types.TopologyDeployRequest `json:"topologyDeployRequest"`
}
