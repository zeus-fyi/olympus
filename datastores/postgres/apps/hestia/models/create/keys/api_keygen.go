package create_keys

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/keys"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"k8s.io/apimachinery/pkg/util/rand"
)

func CreateUserAPIKey(ctx context.Context, ou org_users.OrgUser) (Key, error) {
	k := Key{keys.NewKey()}
	k.PublicKeyName = "API Key"
	k.PublicKey = rand.String(64)
	k.PublicKeyVerified = true
	k.PublicKeyTypeID = keys.BearerKeyTypeID
	k.UserID = ou.UserID
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `WITH keys AS (
				  	DELETE FROM users_keys WHERE user_id = $2 AND public_key_type_id = $5 AND public_key_name = $3
				  )		
				  INSERT INTO users_keys(public_key, user_id, public_key_name, public_key_verified, public_key_type_id)
				  VALUES ($1, $2, $3, $4, $5)
				  `
	r, err := apps.Pg.Exec(ctx, q.RawQuery, k.PublicKey, k.UserID, k.PublicKeyName, k.PublicKeyVerified, k.PublicKeyTypeID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return k, err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("CreateUserAPIKey: %s, Rows Affected: %d", q.LogHeader(Sn), rowsAffected)
	return k, misc.ReturnIfErr(err, q.LogHeader(Sn))
}
