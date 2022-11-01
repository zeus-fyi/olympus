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
	query := `
		WITH cte_insert_top AS (
			INSERT INTO topologies(name) VALUES ($1) RETURNING topology_id
		) INSERT INTO org_users_topologies(topology_id, org_id, user_id)
		  VALUES ((SELECT topology_id FROM cte_insert_top), $2, $3)
		  RETURNING topology_id`
	err := apps.Pg.QueryRowWArgs(ctx, query, t.Name, t.OrgID, t.UserID).Scan(&t.TopologyID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	t.SetTopologyID(t.TopologyID)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
