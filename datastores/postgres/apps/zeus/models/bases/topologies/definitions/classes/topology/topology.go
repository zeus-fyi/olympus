package topology

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/class_type"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
)

type Topology struct {
	autogen_bases.Topologies // -- specific sw names. eg. prysm_validator_client
	class_type.TopologyClassType
	class_type.TopologyClass

	kns.Kns
	state.State
}

func NewTopology() Topology {
	t := Topology{
		Kns: kns.Kns{},
	}
	return t
}

func NewClusterTopology() Topology {
	class := class_type.NewTopologyClass()
	classType := class_type.NewClusterClassTopologyType()
	t := Topology{autogen_bases.Topologies{}, classType, class, kns.Kns{}, state.NewState()}
	return t
}

func (t *Topology) SetTopologyID(id int) {
	t.TopologyID = id
	t.Kns.TopologyID = id
	t.State.TopologyID = id
}

func (t *Topology) SetOrgUserIDs(orgID, userID int) {
	t.State.OrgID = orgID
	t.State.UserID = userID
}
