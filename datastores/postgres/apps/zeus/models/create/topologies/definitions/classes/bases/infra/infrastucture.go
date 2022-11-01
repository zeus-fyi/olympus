package create_infra

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/bases/infra"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/packages"
)

type InfraBaseTopology struct {
	infra.InfraBaseTopology
	packages.Packages
}

func NewCreateInfrastructure() InfraBaseTopology {
	ib := infra.NewInfrastructureBaseTopology()
	pkg := packages.NewPackageInsert()
	ibc := InfraBaseTopology{ib, pkg}
	return ibc
}
