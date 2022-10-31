package create_topology

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (t *Topologies) GetInsertTopologyQueryParams() sql_query_templates.QueryParams {
	q := sql_query_templates.NewQueryParam("InsertTopology", "topologies", "where", 1000, []string{})
	q.TableName = t.Topology.GetTableName()
	q.Columns = t.GetTableColumns()
	q.Values = []apps.RowValues{t.GetRowValues("default")}
	return q
}

func (t *Topologies) InsertTopology(ctx context.Context) error {
	q := t.GetInsertTopologyQueryParams()
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))
	err := apps.Pg.QueryRow(ctx, q.InsertSingleElementQueryReturning("topology_id")).Scan(&t.TopologyID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	t.SetTopologyID(t.TopologyID)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
