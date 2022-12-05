package create_clusters

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/class_types"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/topology/classes/systems"
)

type Cluster struct {
	systems.Systems
}

func NewClusterClassTopologyType(orgID int, name string) Cluster {
	c := Cluster{
		Systems: systems.Systems{TopologySystemComponents: autogen_bases.TopologySystemComponents{
			OrgID:                       orgID,
			TopologyClassTypeID:         class_types.ClusterClassTypeID,
			TopologySystemComponentName: name,
		}},
	}
	return c
}

type TopologyClass struct {
	OrgUser         org_users.OrgUser `json:"orgUser"`
	TopologyClassID int               `json:"topologyClassID"`
	SkeletonBaseIDs []int             `json:"skeletonBaseIDs"`
}
