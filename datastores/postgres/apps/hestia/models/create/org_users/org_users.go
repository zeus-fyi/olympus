package create_org_users

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

type OrgUser struct {
	org_users.OrgUser
}

func NewCreateOrgUser() OrgUser {
	o := OrgUser{org_users.NewOrgUser()}
	return o
}

func NewCreateOrgUserWithOrgID(orgID int) OrgUser {
	o := OrgUser{org_users.NewOrgUser()}
	o.OrgID = orgID
	o.OrgUser.OrgID = orgID
	return o
}
