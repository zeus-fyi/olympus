package create_topology

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/topology"
)

type Topologies struct {
	topology.Topology
}

func NewCreateTopology() Topologies {
	tc := Topologies{topology.NewTopology()}
	return tc
}

func NewCreateOrgUsersInfraTopology(ou org_users.OrgUser) Topologies {
	tc := Topologies{topology.NewOrgUsersInfraTopology(ou)}
	return tc
}

const Sn = "Topologies"
