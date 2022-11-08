package base

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	topology_activities "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities"
)

type TopologyActionRequest struct {
	Action string
	topology_activities.TopologyActivityRequest
}

func CreateTopologyActionRequestWithOrgUser(action string, ou org_users.OrgUser) TopologyActionRequest {
	tar := TopologyActionRequest{
		Action: action,
	}
	tar.OrgUser = ou
	return tar
}
