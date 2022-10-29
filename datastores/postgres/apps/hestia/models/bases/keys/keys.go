package keys

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
)

type KeyGroup struct {
	autogen_bases.UsersKeyGroups
	Keys []Key
}
