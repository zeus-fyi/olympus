package class_type

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

const InfrastructureBaseClassTypeID = 0
const ConfigurationBaseClassTypeID = 1
const SkeletonBaseClassTypeID = 2
const BaseClassTypeID = 3
const ClusterClassTypeID = 4
const MatrixClassTypeID = 5
const SystemClassTypeID = 6

// TopologyClass types skeleton, infrastructure, configuration, base, cluster, matrix, system
type TopologyClass struct {
	Type       autogen_bases.TopologyClassTypes // should be either a base or base subtype, config, cluster, or matrix
	Definition autogen_bases.TopologyClasses
}

func NewTopologyClass() TopologyClass {
	tc := TopologyClass{
		Type:       autogen_bases.TopologyClassTypes{},
		Definition: autogen_bases.TopologyClasses{},
	}
	return tc
}

func (t *TopologyClass) GetClassDefinition() autogen_bases.TopologyClasses {
	return t.Definition
}

func (t *TopologyClass) SetClassDefinition(def autogen_bases.TopologyClasses) {
	t.Definition = def
}
