package topology

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/class_type"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
)

type Topology struct {
	autogen_bases.Topologies // -- specific sw names. eg. prysm_validator_client
	class_type.TopologyClassType
	class_type.TopologyClass
	OrgUserTopology

	kns.TopologyKubeCtxNs
	topology_deployment_status.Status
}

func NewTopology() Topology {
	t := Topology{
		TopologyKubeCtxNs: kns.TopologyKubeCtxNs{},
	}
	return t
}

func NewClusterTopology() Topology {
	class := class_type.NewTopologyClass()
	classType := class_type.NewClusterClassTopologyType()
	out := NewOrgUserTopology()
	status := topology_deployment_status.NewTopologyStatus()
	t := Topology{autogen_bases.Topologies{}, classType, class, out, kns.TopologyKubeCtxNs{}, status}
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

func NewOrgUsersInfraTopology(ou org_users.OrgUser) Topology {
	class := class_type.NewTopologyClass()
	classType := class_type.NewInfraClassTopologyType()
	out := NewOrgUserTopologyFromOrgUser(ou)
	status := topology_deployment_status.NewTopologyStatus()
	t := Topology{autogen_bases.Topologies{}, classType, class, out, kns.TopologyKubeCtxNs{}, status}
	return t
}
