package infra

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
)

type InfraBaseTopology struct {
	Name string
	org_users.OrgUser
	autogen_bases.TopologyInfrastructureComponents
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

func NewOrgUserInfrastructureBaseTopology(topologyName string, ou org_users.OrgUser) InfraBaseTopology {
	tic := autogen_bases.TopologyInfrastructureComponents{
		TopologyID:     0,
		ChartPackageID: 0,
	}
	inf := InfraBaseTopology{topologyName, ou, tic}
	return inf
}

func (i *InfraBaseTopology) AddChartPackageID(id int) {
	i.ChartPackageID = id
}
