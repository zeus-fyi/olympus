package auth

import (
	"context"
	"errors"
	"time"

	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/keys"
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
)

func VerifyBearerToken(ctx context.Context, token string) (read_keys.OrgUserKey, error) {
	key := read_keys.OrgUserKey{
		Key: keys.Key{
			UsersKeys: autogen_bases.UsersKeys{
				UserID:            0,
				PublicKeyName:     "",
				PublicKeyVerified: false,
				PublicKeyTypeID:   keys.BearerKeyTypeID,
				CreatedAt:         time.Time{},
				PublicKey:         token,
			},
			KeyType: keys.NewBearerKeyType(),
		},
	}
	err := key.VerifyUserBearerToken(ctx)
	if err != nil {
		return read_keys.OrgUserKey{}, err
	}
	if key.PublicKeyVerified == false {
		return read_keys.OrgUserKey{}, errors.New("unauthorized key")
	}
	return key, err
}

func VerifyBearerTokenService(ctx context.Context, token, serviceName string) (read_keys.OrgUserKey, error) {
	key := read_keys.OrgUserKey{
		Key: keys.Key{
			UsersKeys: autogen_bases.UsersKeys{
				UserID:            0,
				PublicKeyName:     "",
				PublicKeyVerified: false,
				CreatedAt:         time.Time{},
				PublicKey:         token,
			},
			KeyType: keys.NewBearerKeyType(),
		},
	}
	err := key.VerifyUserTokenService(ctx, serviceName)
	if err != nil {
		return read_keys.OrgUserKey{}, err
	}
	if key.PublicKeyVerified == false {
		return read_keys.OrgUserKey{}, errors.New("unauthorized key")
	}
	return key, err
}
