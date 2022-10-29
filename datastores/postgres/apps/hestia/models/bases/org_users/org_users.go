package org_users

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/orgs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/users"
)

type OrgUser struct {
	orgs.Org
	users.User
}

type OrgUsers struct {
	orgs.Org
	users.User
}
