package auth

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/keys"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
)

const (
	TemporalOrgID  = 7138983863666903883
	TemporalUserID = 7138958574876245567
	AdminUserID    = 7138958574876245565
)

func VerifyInternalBearerToken(ctx context.Context, token string) (read_keys.OrgUserKey, error) {
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
		log.Info().Int("userID", key.UserID).Int("orgID", key.OrgID).Err(err)
		return read_keys.OrgUserKey{}, err
	}
	if key.PublicKeyVerified == false {
		log.Info().Int("userID", key.UserID).Int("orgID", key.OrgID)
		return read_keys.OrgUserKey{}, errors.New("unauthorized key")
	}

	if key.GetUserID() != TemporalUserID && key.OrgID != TemporalOrgID {
		log.Info().Int("userID", key.UserID).Int("orgID", key.OrgID)
		return read_keys.OrgUserKey{}, errors.New("unauthorized key")
	}
	return key, err
}

const (
	SamsOrgID = 1701381301753642000
)

func VerifyInternalAdminBearerToken(ctx context.Context, token string) (read_keys.OrgUserKey, error) {
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
		log.Info().Int("userID", key.UserID).Int("orgID", key.OrgID).Err(err)
		return read_keys.OrgUserKey{}, err
	}
	if key.PublicKeyVerified == false {
		log.Info().Int("userID", key.UserID).Int("orgID", key.OrgID)
		return read_keys.OrgUserKey{}, errors.New("unauthorized key")
	}
	if key.OrgID == SamsOrgID {
		return key, nil
	}
	if key.GetUserID() != AdminUserID && key.OrgID != TemporalOrgID {
		log.Info().Int("userID", key.UserID).Int("orgID", key.OrgID)
		return read_keys.OrgUserKey{}, errors.New("unauthorized key")
	}
	return key, err
}

func FetchTemporalAuthToken(ctx context.Context) (read_keys.OrgUserKey, error) {
	key := read_keys.OrgUserKey{}
	ou := org_users.NewOrgUserWithID(TemporalOrgID, TemporalUserID)
	err := key.QueryUserBearerToken(ctx, ou)
	if err != nil {
		log.Err(err).Msg("FetchTemporalAuthToken, failed to query for auth token")
		return read_keys.OrgUserKey{}, err
	}
	return key, err
}

func FetchUserAuthToken(ctx context.Context, ou org_users.OrgUser) (read_keys.OrgUserKey, error) {
	key := read_keys.OrgUserKey{}
	err := key.QueryUserBearerToken(ctx, ou)
	if err != nil {
		log.Err(err).Msg("FetchUserAuthToken, failed to query for auth token")
		return read_keys.OrgUserKey{}, err
	}
	return key, err
}

func FetchUserAuthTokenDiscord(ctx context.Context, userId int) (string, error) {
	kv, err := read_keys.GetDiscordKey(ctx, userId)
	if err != nil {
		log.Err(err).Msg("FetchUserAuthToken, failed to query for auth token")
		return "", err
	}
	return kv, err

}
