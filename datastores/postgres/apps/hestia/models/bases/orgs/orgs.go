package orgs

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"

type Org struct {
	autogen_bases.Orgs
}

func NewOrg() Org {
	o := Org{autogen_bases.Orgs{
		OrgID:    0,
		Metadata: "{}",
	}}
	return o
}
