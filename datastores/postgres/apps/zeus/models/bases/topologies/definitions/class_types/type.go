package class_types

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

const InfrastructureBaseClassTypeID = 0
const ConfigurationBaseClassTypeID = 1
const SkeletonBaseClassTypeID = 2
const BaseClassTypeID = 3
const ClusterClassTypeID = 4
const MatrixClassTypeID = 5
const SystemClassTypeID = 6

type TopologyClassType struct {
	autogen_bases.TopologyClassTypes // should be either a base or base subtype, config, cluster, or matrix
}

func NewInfraClassTopologyType() TopologyClassType {
	tc := TopologyClassType{
		autogen_bases.TopologyClassTypes{
			TopologyClassTypeID: InfrastructureBaseClassTypeID,
		},
	}
	return tc
}
