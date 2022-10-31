package create_topology

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/topology"
)

type Topologies struct {
	topology.Topology
}

func NewCreateTopology() Topologies {
	tc := Topologies{topology.NewTopology()}
	return tc
}

const Sn = "Topologies"
