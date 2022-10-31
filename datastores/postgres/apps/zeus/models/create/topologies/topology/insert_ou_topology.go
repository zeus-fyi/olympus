package create_topology

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (t *Topologies) GetInsertOrgUsersTopologyQueryParams() sql_query_templates.QueryParams {
	q := sql_query_templates.NewQueryParam("InsertOrgUsersTopology", "org_users_topologies", "where", 1000, []string{})
	q.TableName = t.OrgUsersTopologies.GetTableName()
	q.Columns = t.OrgUsersTopologies.GetTableColumns()
	q.Values = []apps.RowValues{t.OrgUsersTopologies.GetRowValues("default")}
	return q
}

func (t *Topologies) InsertOrgUsersTopology(ctx context.Context) error {
	q := t.GetInsertOrgUsersTopologyQueryParams()
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))

	query := q.InsertSingleElementQuery()
	r, err := apps.Pg.Exec(ctx, query)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("OrgUsersTopology: %s, Rows Affected: %d", q.LogHeader(Sn), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
