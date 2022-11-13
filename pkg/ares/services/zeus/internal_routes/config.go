package internal_routes

import (
	"time"

	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
)

var status = topology_deployment_status.Status{
	TopologyKubeCtxNs: kns.TopologyKubeCtxNs{TopologiesKns: autogen_bases.TopologiesKns{
		TopologyID:    1668293626655515904,
		CloudProvider: "do",
		Region:        "sfo3",
		Context:       "dev-sfo3-zeus",
		Namespace:     "demo",
		Env:           "dev",
	}},
	DeployStatus: topology_deployment_status.DeployStatus{TopologiesDeployed: autogen_bases.TopologiesDeployed{
		DeploymentID:   0,
		TopologyID:     1668293626655515904,
		TopologyStatus: "Rollback",
		UpdatedAt:      time.Time{},
	}},
}
