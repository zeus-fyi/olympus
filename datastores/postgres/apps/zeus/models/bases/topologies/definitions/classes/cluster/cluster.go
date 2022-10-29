package clusters

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/bases"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/topology"
)

type Cluster struct {
	topology.Topology
	Base       bases.Base
	Components topology.Dependencies
}

func NewCluster() Cluster {
	cl := Cluster{topology.NewClusterTopology(), bases.NewBase(), []topology.Component{}}
	return cl
}
