package zeus_demo

import (
	create_or_update_deploy "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/external/create_or_update"
	destroy_deploy_request "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/external/destroy"
)

var obolDeployKnsReq = create_or_update_deploy.TopologyDeployRequest{
	TopologyID:    1668723930348846000,
	CloudProvider: "do",
	Region:        "sfo3",
	Context:       "do-sfo3-dev-do-sfo3-zeus",
	Namespace:     "obol",
	Env:           "dev",
}

var obolDestroyDeployKnsReq = destroy_deploy_request.TopologyDestroyDeployRequest{
	TopologyID:    1668723930348846000,
	CloudProvider: "do",
	Region:        "sfo3",
	Context:       "do-sfo3-dev-do-sfo3-zeus",
	Namespace:     "obol",
	Env:           "dev",
}

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
