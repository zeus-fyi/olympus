package clusters

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
)

type Cluster struct {
	classes.TopologyDependentComponent
}

func NewCluster() Cluster {
	cl := Cluster{classes.TopologyDependentComponent{
		TopologyDependentComponents:      autogen_bases.TopologyDependentComponents{},
		TopologyInfrastructureComponents: autogen_bases.TopologyInfrastructureComponents{},
		InfraComponents:                  nil,
		Status:                           state.State{},
	}}
	return cl
}
