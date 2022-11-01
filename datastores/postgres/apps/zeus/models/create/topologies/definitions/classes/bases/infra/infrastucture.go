package create_infra

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/bases/infra"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/packages"
)

type InfraBaseTopology struct {
	infra.InfraBaseTopology

	packages.Packages
}

func NewCreateInfrastructure() InfraBaseTopology {
	pkg := packages.NewPackageInsert()
	infc := infra.NewInfrastructureBaseTopology()
	ibc := InfraBaseTopology{infc, pkg}
	return ibc
}

func NewOrgUserCreateInfraFromNativeK8s(topologyName string, ou org_users.OrgUser, pkg packages.Packages) InfraBaseTopology {
	infc := infra.NewOrgUserInfrastructureBaseTopology(topologyName, ou)
	ibc := InfraBaseTopology{infc, pkg}
	return ibc
}
