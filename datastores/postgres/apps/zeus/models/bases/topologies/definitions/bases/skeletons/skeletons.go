package skeletons

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

type Skeleton struct {
	autogen_bases.TopologySkeletonBaseComponents
}

type SkeletonBase struct {
	autogen_bases.TopologyInfrastructureComponents
}
