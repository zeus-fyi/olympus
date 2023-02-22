package read_keys

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/keys"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

const Sn = "KeyReader"

func (k *OrgUserKey) QueryVerifyUserBearerToken() sql_query_templates.QueryParams {
	var q sql_query_templates.QueryParams
	query := fmt.Sprintf(`
	SELECT usk.public_key_verified, ou.org_id, ou.user_id
	FROM users_keys usk
	INNER JOIN key_types kt ON kt.key_type_id = usk.public_key_type_id
	INNER JOIN org_users ou ON ou.user_id = usk.user_id
	WHERE public_key = $1
	`)
	q.RawQuery = query
	return q
}
func (k *OrgUserKey) VerifyUserBearerToken(ctx context.Context) error {
	q := k.QueryVerifyUserBearerToken()
	log.Debug().Interface("VerifyUserBearerToken:", q.LogHeader(Sn))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, k.PublicKey).Scan(&k.PublicKeyVerified, &k.OrgID, &k.UserID)
	if err != nil {
		k.PublicKeyVerified = false
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func (k *OrgUserKey) QueryVerifyUserBearerTokenService() sql_query_templates.QueryParams {
	var q sql_query_templates.QueryParams
	query := fmt.Sprintf(`
	SELECT usk.public_key_verified, ou.org_id, ou.user_id
	FROM users_keys usk
	INNER JOIN key_types kt ON kt.key_type_id = usk.public_key_type_id
	INNER JOIN users_key_services uksvc ON uksvc.public_key = usk.public_key
	INNER JOIN services svcs ON svcs.service_id = uksvc.service_id
	INNER JOIN org_users ou ON ou.user_id = usk.user_id
	WHERE usk.public_key = $1 AND service_name = $2
	`)
	q.RawQuery = query
	return q
}
func (k *OrgUserKey) VerifyUserBearerTokenService(ctx context.Context, serviceName string) error {
	q := k.QueryVerifyUserBearerTokenService()
	log.Debug().Interface("QueryVerifyUserBearerTokenService:", q.LogHeader(Sn))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, k.PublicKey, serviceName).Scan(&k.PublicKeyVerified, &k.OrgID, &k.UserID)
	if err != nil {
		k.PublicKeyVerified = false
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func (k *OrgUserKey) GetUserID() int {
	return k.UserID
}

func (k *OrgUserKey) QueryUserToken() sql_query_templates.QueryParams {
	var q sql_query_templates.QueryParams
	query := fmt.Sprintf(`
	SELECT usk.public_key
	FROM users_keys usk
	INNER JOIN key_types kt ON kt.key_type_id = usk.public_key_type_id
	INNER JOIN org_users ou ON ou.user_id = usk.user_id
	WHERE ou.org_id = $1 AND usk.user_id = $2 AND usk.public_key_type_id = $3
	`)
	q.RawQuery = query
	return q
}

func (k *OrgUserKey) QueryUserBearerToken(ctx context.Context, ou org_users.OrgUser) error {
	q := k.QueryUserToken()
	log.Debug().Interface("QueryUserBearerToken:", q.LogHeader(Sn))
	k.OrgID = ou.OrgID
	k.UserID = ou.UserID
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, ou.OrgID, ou.UserID, keys.BearerKeyTypeID).Scan(&k.PublicKey)
	if err != nil {
		k.PublicKeyVerified = false
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
