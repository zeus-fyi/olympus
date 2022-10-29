package class_type

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

type TopologyClass struct {
	autogen_bases.TopologyClassTypes
	autogen_bases.TopologyClasses
}

func NewTopologyClass() TopologyClass {
	tc := TopologyClass{
		TopologyClassTypes: autogen_bases.TopologyClassTypes{},
		TopologyClasses:    autogen_bases.TopologyClasses{},
	}
	return tc
}
