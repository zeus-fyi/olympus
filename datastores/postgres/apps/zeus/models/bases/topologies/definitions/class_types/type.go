package class_types

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
	autogen_bases.TopologyClasses
}
type TopologyClassType struct {
	autogen_bases.TopologyClassTypes // should be either a base or base subtype, config, cluster, or matrix
}

func NewTopologyClass() TopologyClass {
	tc := TopologyClass{
		autogen_bases.TopologyClasses{
			TopologyClassID:     0,
			TopologyClassTypeID: 0,
			TopologyClassName:   "",
		},
	}
	return tc
}

func NewInfraClassTopologyType() TopologyClassType {
	tc := TopologyClassType{
		autogen_bases.TopologyClassTypes{
			TopologyClassTypeID: InfrastructureBaseClassTypeID,
		},
	}
	return tc
}
