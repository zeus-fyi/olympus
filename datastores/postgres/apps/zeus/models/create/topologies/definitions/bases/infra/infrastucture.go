package create_infra

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/bases/infra"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/packages"
)

type InfraBaseTopology struct {
	infra.InfraBaseTopology

	packages.Packages

	SkeletonBaseID  int `json:"skeletonBaseID,omitempty"`
	TopologyClassID int `json:"topologyClassID,omitempty"`
}

func NewCreateInfrastructure() InfraBaseTopology {
	pkg := packages.NewPackageInsert()
	infc := infra.NewInfrastructureBaseTopology()
	ibc := InfraBaseTopology{infc, pkg, 0, 0}
	return ibc
}
