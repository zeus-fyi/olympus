package topology

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
)

type OrgUserTopology struct {
	autogen_bases.OrgUsersTopologies
}

func NewOrgUserTopology() OrgUserTopology {
	ou := org_users.NewOrgUser()
	out := NewOrgUserTopologyFromOrgUser(ou)
	return out
}

func NewOrgUserTopologyFromOrgUser(ou org_users.OrgUser) OrgUserTopology {
	out := OrgUserTopology{autogen_bases.OrgUsersTopologies{
		TopologyID: 0,
		OrgID:      ou.OrgID,
		UserID:     ou.UserID,
	}}
	return out
}
