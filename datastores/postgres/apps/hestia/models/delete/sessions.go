package hestia_delete

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/keys"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

const Sn = "delete_sessions"

func DeleteUserSessionKey(ctx context.Context, sessionToken string) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `DELETE FROM users_keys WHERE (public_key = $1 AND public_key_type_id = $2)`
	r, err := apps.Pg.Exec(ctx, q.RawQuery, sessionToken, keys.SessionIDKeyTypeID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("DeleteUserSessionKey: %s, Rows Affected: %d", q.LogHeader(Sn), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
