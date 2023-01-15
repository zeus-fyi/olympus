package create_orgs

import "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/orgs"

type Org struct {
	orgs.Org
}

func NewCreateOrg() Org {
	o := Org{orgs.NewOrg()}
	return o
}

func NewCreateNamedOrg(orgName string) Org {
	o := Org{orgs.NewOrg()}
	o.Name = orgName
	return o
}
