package skeletons

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/class_types"
)

type SkeletonBase struct {
	autogen_bases.TopologySkeletonBaseComponents
}

func NewSkeletonBase(orgID, baseComponentID int, name string) SkeletonBase {
	s := SkeletonBase{autogen_bases.TopologySkeletonBaseComponents{
		OrgID:                    orgID,
		TopologyBaseComponentID:  baseComponentID,
		TopologyClassTypeID:      class_types.SkeletonBaseClassTypeID,
		TopologySkeletonBaseName: name,
	}}
	return s
}
