package org_users

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
)

type OrgUser struct {
	autogen_bases.OrgUsers `json:"orgUser"`
}

func NewOrgUser() OrgUser {
	o := OrgUser{autogen_bases.OrgUsers{
		OrgID:  0,
		UserID: 0,
	}}
	return o
}

func NewOrgUserWithID(orgID, userID int) OrgUser {
	o := OrgUser{autogen_bases.OrgUsers{
		OrgID:  orgID,
		UserID: userID,
	}}
	return o
}
