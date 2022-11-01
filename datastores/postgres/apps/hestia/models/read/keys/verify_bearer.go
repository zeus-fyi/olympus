package read_keys

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
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

	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, k.PublicKey).Scan(&k.PublicKeyVerified, &k.OrgID, &k.OrgUser.UserID)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
