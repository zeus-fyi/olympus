package topology

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/class_type"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
)

type Topology struct {
	Name string
	kns.TopologyKubeCtxNs
	class_type.TopologyClassType
	class_type.TopologyClass
	OrgUserTopology
}

func NewTopology() Topology {
	t := Topology{
		TopologyKubeCtxNs: kns.TopologyKubeCtxNs{},
	}
	return t
}

func (t *Topology) SetTopologyID(id int) {
	t.TopologyKubeCtxNs.TopologyID = id
	t.OrgUserTopology.TopologyID = id
}

func (t *Topology) SetOrgUserIDs(orgID, userID int) {
	t.OrgID = orgID
	t.UserID = userID
}
