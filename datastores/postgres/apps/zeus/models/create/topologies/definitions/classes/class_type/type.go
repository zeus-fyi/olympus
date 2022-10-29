package create_class_type

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/class_type"
)

type TopologyClass struct {
	class_type.TopologyClass
}

func NewCreateTopologyClass() TopologyClass {
	ct := class_type.NewTopologyClass()
	class := TopologyClass{ct}
	return class
}
