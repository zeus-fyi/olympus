package create_class_type

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/class_types"
)

type TopologyClass struct {
	class_types.TopologyClass
}

func NewCreateTopologyClass() TopologyClass {
	ct := class_types.NewTopologyClass()
	class := TopologyClass{ct}
	return class
}
