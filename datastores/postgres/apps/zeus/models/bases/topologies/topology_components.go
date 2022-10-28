package topologies

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

type TopologyDependentComponent struct {
	autogen_bases.TopologyDependentComponents
	autogen_bases.TopologyInfrastructureComponents
	// use map->int to represent hierarchy ordering
	// all same level uncoupled items are on the same depth
	InfraComponents map[int][]TopologyDependentComponent
	Status          autogen_bases.TopologiesDeployed
}
