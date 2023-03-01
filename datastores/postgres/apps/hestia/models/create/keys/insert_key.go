package create_keys

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

const Sn = "UserKey"

func (k *Key) InsertUserKey(ctx context.Context) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
				  INSERT INTO users_keys(public_key, user_id, public_key_name, public_key_verified, public_key_type_id)
				  VALUES ($1, $2, $3, $4, $5)
				  `
	r, err := apps.Pg.Exec(ctx, q.RawQuery, k.PublicKey, k.UserID, k.PublicKeyName, k.PublicKeyVerified, k.PublicKeyTypeID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("OrgUser: %s, Rows Affected: %d", q.LogHeader(Sn), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
