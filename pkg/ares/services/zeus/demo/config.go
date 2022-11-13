package zeus_demo

import (
	create_or_update_deploy "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/external/create_or_update"
	destroy_deploy_request "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/external/destroy"
)

var deployKnsReq = create_or_update_deploy.TopologyDeployRequest{
	TopologyID:    1668320728359007000,
	CloudProvider: "do",
	Region:        "sfo3",
	Context:       "dev-sfo3-zeus",
	Namespace:     "demo",
	Env:           "dev",
}

var deployDestroyKnsReq = destroy_deploy_request.TopologyDestroyDeployRequest{
	TopologyID:    1668320728359007000,
	CloudProvider: "do",
	Region:        "sfo3",
	Context:       "dev-sfo3-zeus",
	Namespace:     "demo",
	Env:           "dev",
}
