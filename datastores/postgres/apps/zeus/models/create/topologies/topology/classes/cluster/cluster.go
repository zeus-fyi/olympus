package create_clusters

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/bases"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/class_types"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/topology/classes/systems"
)

type Cluster struct {
	systems.Systems
	Bases []bases.Base `json:"bases,omitempty"`
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

func NewClusterClassTopologyTypeWithBases(orgID int, name string, baes []string) Cluster {
	basesToInsert := []bases.Base{}
	for _, v := range baes {
		b := bases.NewBaseClassTopologyType(orgID, 0, v)
		basesToInsert = append(basesToInsert, b)
	}
	c := Cluster{
		Systems: systems.Systems{TopologySystemComponents: autogen_bases.TopologySystemComponents{
			OrgID:                       orgID,
			TopologyClassTypeID:         class_types.ClusterClassTypeID,
			TopologySystemComponentName: name,
		}},
		Bases: basesToInsert,
	}

	return c
}
