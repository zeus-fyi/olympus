package create_keys

import "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/keys"

type Key struct {
	keys.Key
}

func NewCreateKey(userID int, pubkey string) Key {
	k := keys.NewKeyForUser(userID, pubkey)
	ck := Key{k}
	return ck
}
