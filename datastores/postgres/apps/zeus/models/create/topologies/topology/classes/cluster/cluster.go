package create_clusters

import (
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
