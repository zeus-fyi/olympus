package hestia_compute_resources

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
	cte.SubCTEs = []sql_query_templates.SubCTE{}
	cte.Params = []interface{}{}

	for _, node := range nodes {
		tmp := &node
		re := autogen_bases.Resources{
			ResourceID: ts.UnixTimeStampNow(),
			Type:       "node",
		}
		node.ResourceID = re.ResourceID
		queryResourceId := fmt.Sprintf("resource_id_insert_%d", node.ResourceID)
		scteRe := sql_query_templates.NewSubInsertCTE(queryResourceId)
		scteRe.TableName = re.GetTableName()
		scteRe.Columns = re.GetTableColumns()
		scteRe.Values = []apps.RowValues{re.GetRowValues(queryResourceId)}
		cte.SubCTEs = append(cte.SubCTEs, scteRe)
		queryName := fmt.Sprintf("nodes_insert_%d", node.ResourceID)
		scte := sql_query_templates.NewSubInsertCTE(queryName)
		scte.TableName = tmp.GetTableName()
		scte.Columns = tmp.GetTableColumns()
		scte.Values = []apps.RowValues{tmp.GetRowValues(queryName)}
		cte.SubCTEs = append(cte.SubCTEs, scte)
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
