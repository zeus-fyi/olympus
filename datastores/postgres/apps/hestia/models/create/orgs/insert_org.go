package create_orgs

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

const Sn = "Org"

func (o *Org) InsertOrg(ctx context.Context) error {
	q := sql_query_templates.NewQueryParam("InsertOrg", "orgs", "where", 1000, []string{})
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))

	q.RawQuery = `INSERT INTO orgs (org_id, name, metadata) VALUES ($1, $2, $3)`

	r, err := apps.Pg.Exec(ctx, q.RawQuery, o.OrgID, o.Name, o.Metadata)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("InsertOrg: %s, Rows Affected: %d", q.LogHeader(Sn), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
