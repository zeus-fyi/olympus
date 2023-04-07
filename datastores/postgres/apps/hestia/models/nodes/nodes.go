package hestia_nodes

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

const Sn = "Nodes"

var ts chronos.Chronos

func InsertNodes(ctx context.Context, nodes autogen_bases.NodesSlice) error {
	q := sql_query_templates.NewQueryParam("InsertNodes", "nodes", "where", 1000, []string{})
	cte := sql_query_templates.CTE{Name: "InsertNodes"}
	cte.SubCTEs = make([]sql_query_templates.SubCTE, len(nodes))
	cte.Params = []interface{}{}

	for i, node := range nodes {
		tmp := &node
		node.NodeID = ts.UnixTimeStampNow()
		queryName := fmt.Sprintf("nodes_insert_%d", node.NodeID)
		scte := sql_query_templates.NewSubInsertCTE(queryName)
		scte.TableName = tmp.GetTableName()
		scte.Columns = tmp.GetTableColumns()
		scte.Values = []apps.RowValues{tmp.GetRowValues(queryName)}
		cte.SubCTEs[i] = scte
	}
	q.RawQuery = cte.GenerateChainedCTE()
	r, err := apps.Pg.Exec(ctx, q.RawQuery, cte.Params...)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("Nodes: %s, Rows Affected: %d", q.LogHeader(Sn), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
