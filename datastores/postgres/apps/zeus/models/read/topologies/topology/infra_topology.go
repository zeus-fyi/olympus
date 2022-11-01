package read_topology

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/bases/infra"
	read_charts "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/charts"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type InfraBaseTopology struct {
	infra.InfraBaseTopology
	read_charts.Chart
}

func NewInfraTopologyReader() InfraBaseTopology {
	bt := infra.NewInfrastructureBaseTopology()
	cr := read_charts.NewChartReader()
	rt := InfraBaseTopology{bt, cr}
	return rt
}

const Sn = "ReadInfraBaseTopology"

func (t *InfraBaseTopology) SelectInfraTopologyQuery() sql_query_templates.QueryParams {
	var q sql_query_templates.QueryParams
	q.QueryName = "SelectInfraTopologyQuery"
	q.CTEQuery.Params = append(q.CTEQuery.Params, t.TopologyID, t.OrgID, t.UserID)
	q.RawQuery = read_charts.FetchChartQuery(q)
	return q
}

func (t *InfraBaseTopology) SelectTopology(ctx context.Context) error {
	q := t.SelectInfraTopologyQuery()

	log.Debug().Interface("SelectTopologyQuery", q.LogHeader(Sn))
	err := t.SelectSingleChartsResources(ctx, q)
	return err
}
