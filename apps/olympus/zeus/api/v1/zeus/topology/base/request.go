package base

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

type TopologyActionRequest struct {
	Action string

	org_users.OrgUser
}

func CreateTopologyActionRequestWithOrgUser(action string, ou org_users.OrgUser) TopologyActionRequest {
	tar := TopologyActionRequest{
		Action:  action,
		OrgUser: ou,
	}
	return tar
}
