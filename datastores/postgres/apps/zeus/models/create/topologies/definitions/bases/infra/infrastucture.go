package create_infra

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/bases/infra"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/packages"
)

type InfraBaseTopology struct {
	infra.InfraBaseTopology

	packages.Packages

	Tag              string `json:"tag,omitempty"`
	SkeletonBaseName string `json:"skeletonBaseName,omitempty"`
	TopologyClassID  int    `json:"topologyClassID,omitempty"`
}

func NewCreateInfrastructure() InfraBaseTopology {
	pkg := packages.NewPackageInsert()
	infc := infra.NewInfrastructureBaseTopology()
	ibc := InfraBaseTopology{infc, pkg, "latest-internal", "", 0}
	return ibc
}
