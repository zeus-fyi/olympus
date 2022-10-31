package base

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

type TopologyActionRequest struct {
	Action string

	org_users.OrgUser
}
