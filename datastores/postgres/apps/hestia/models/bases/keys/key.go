package keys

import (
	"time"

	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
)

type Key struct {
	autogen_bases.UsersKeys
	KeyType
}

func NewKeyForUser(userID int, publicKey string) Key {
	k := Key{
		UsersKeys: autogen_bases.UsersKeys{
			UserID:            userID,
			PublicKeyName:     "",
			PublicKeyVerified: false,
			PublicKeyTypeID:   0,
			CreatedAt:         time.Time{},
			PublicKey:         publicKey,
		},
	}
	return k
}
