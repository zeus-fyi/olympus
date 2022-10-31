package topology

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/class_type"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
)

type InfraTopology struct {
	Topology
	ChartID int
}

func NewInfraTopology() Topology {
	class := class_type.NewTopologyClass()
	classType := class_type.NewInfraClassTopologyType()
	ou := NewOrgUserTopology()
	t := Topology{autogen_bases.Topologies{}, classType, class, ou, kns.Kns{}, state.NewState()}
	return t
}
