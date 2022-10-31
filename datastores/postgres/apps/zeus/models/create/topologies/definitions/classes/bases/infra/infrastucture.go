package create_infra

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/bases/infra"
)

type InfraBaseTopology struct {
	infra.InfraBaseTopology
}

func NewCreateInfrastructure() InfraBaseTopology {
	ib := infra.NewInfrastructureBaseTopology()
	ibc := InfraBaseTopology{ib}
	return ibc
}
