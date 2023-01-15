package create_org_users

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

const Sn = "OrgUser"

func (o *OrgUser) InsertOrgUser(ctx context.Context, metadata []byte) error {
	q := sql_query_templates.NewQueryParam("NewTestOrgUser", "org_users", "where", 1000, []string{})
	q.RawQuery = `WITH new_user_id AS (
					INSERT INTO users(metadata)
					VALUES ($2)
					RETURNING user_id
				)
					INSERT INTO org_users(org_id, user_id)
					VALUES($1, (SELECT user_id FROM new_user_id))
					RETURNING (SELECT user_id FROM new_user_id)
				`
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))
	var userID int64
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, o.OrgID, metadata).Scan(&userID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	o.UserID = int(userID)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
