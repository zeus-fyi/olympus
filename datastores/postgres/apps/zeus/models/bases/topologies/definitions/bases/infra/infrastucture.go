package infra

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
)

type InfraBaseTopology struct {
	Name                                           string `json:"name"`
	org_users.OrgUser                              `json:"orgUser"`
	autogen_bases.TopologyInfrastructureComponents `json:"topologyInfrastructureComponents"`
}

func NewInfrastructureBaseTopologyWithOrgUser(ou org_users.OrgUser) InfraBaseTopology {
	tic := autogen_bases.TopologyInfrastructureComponents{
		TopologyID:     0,
		ChartPackageID: 0,
	}
	inf := InfraBaseTopology{"", ou, tic}
	return inf
}

func NewInfrastructureBaseTopology() InfraBaseTopology {
	tic := autogen_bases.TopologyInfrastructureComponents{
		TopologyID:     0,
		ChartPackageID: 0,
	}
	ou := org_users.NewOrgUser()
	inf := InfraBaseTopology{"", ou, tic}
	return inf
}

func (i *InfraBaseTopology) AddChartPackageID(id int) {
	i.ChartPackageID = id
}
