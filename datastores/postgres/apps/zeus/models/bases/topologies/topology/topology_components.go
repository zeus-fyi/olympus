package topology

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/bases"
)

// Component single component, can link components to build higher level topologies eg. beacon + exec = full eth2 beacon cluster
type Component struct {
	bases.Base
	autogen_bases.TopologyDependentComponents
}

type Dependencies []Component

func NewTopologyComponent() Component {
	tc := Component{bases.NewBase(),
		autogen_bases.TopologyDependentComponents{
			TopologyClassID: 0,
			TopologyID:      0,
		},
	}
	return tc
}
