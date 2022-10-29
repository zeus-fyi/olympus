package classes

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
)

type TopologyDependentComponent struct {
	// -- links components to build higher level topologies eg. beacon + exec = full eth2 beacon cluster
	autogen_bases.TopologyDependentComponents
	autogen_bases.TopologyInfrastructureComponents
	// use map->int to represent hierarchy ordering
	// all same level uncoupled items are on the same depth
	InfraComponents map[int][]charts.Chart
	Status          state.State
}
