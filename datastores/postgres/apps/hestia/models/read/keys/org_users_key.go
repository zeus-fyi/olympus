package read_keys

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/keys"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

type OrgUserKey struct {
	keys.Key
	org_users.OrgUser
}

func NewKeyReader() OrgUserKey {
	k := keys.NewKey()
	ou := org_users.NewOrgUser()
	return OrgUserKey{k, ou}
}
