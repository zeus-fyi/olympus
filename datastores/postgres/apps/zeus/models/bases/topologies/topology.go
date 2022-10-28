package topologies

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/kns"
)

type Topology struct {
	kns.Kns
	TopologyClass
	// use map->int to represent hierarchy
	// all same level uncoupled items are on the same depth
	TopologyMatrix map[int][]TopologyDependentComponent
}
