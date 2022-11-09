package read_keys

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/keys"
)

type OrgUserKey struct {
	OrgID int
	keys.Key
}

func NewKeyReader() OrgUserKey {
	return OrgUserKey{}
}
