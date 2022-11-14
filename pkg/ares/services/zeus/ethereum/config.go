package zeus_ethereum

import (
	create_or_update_deploy "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/external/create_or_update"
	destroy_deploy_request "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/external/destroy"
)

var deployConsensusClientKnsReq = create_or_update_deploy.TopologyDeployRequest{
	TopologyID:    1668386117729324000,
	CloudProvider: "do",
	Region:        "sfo3",
	Context:       "dev-sfo3-zeus",
	Namespace:     "ethereum",
	Env:           "dev",
}

var deployDestroyConsensusClientKnsReq = destroy_deploy_request.TopologyDestroyDeployRequest{
	TopologyID:    1668386117729324000,
	CloudProvider: "do",
	Region:        "sfo3",
	Context:       "dev-sfo3-zeus",
	Namespace:     "ethereum",
	Env:           "dev",
}
