package create_topology

import (
	"context"
	"fmt"

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
	query := fmt.Sprintf(`
		WITH cte_insert_top AS (
			INSERT INTO topologies(name) VALUES ('%s') RETURNING topology_id
		) INSERT INTO org_users_topologies(topology_id, org_id, user_id)
		  VALUES ((SELECT topology_id FROM cte_insert_top), %d, %d)
		  RETURNING topology_id`,
		t.Name, t.OrgID, t.UserID)
	err := apps.Pg.QueryRow(ctx, query).Scan(&t.TopologyID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	t.SetTopologyID(t.TopologyID)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
