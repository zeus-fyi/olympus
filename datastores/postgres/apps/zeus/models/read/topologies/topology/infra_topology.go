package read_topology

import (
	"context"
	"fmt"

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
	query := fmt.Sprintf(`(
			SELECT tic.chart_package_id
			FROM topology_infrastructure_components tic
			INNER JOIN topologies top ON top.topology_id = tic.topology_id
			INNER JOIN org_users_topologies out ON out.topology_id = tic.topology_id
			WHERE tic.topology_id = %d AND org_id = %d AND user_id = %d
			LIMIT 1
			)
	`, t.TopologyID, t.OrgID, t.UserID)
	q.RawQuery = read_charts.FetchChartQuery(query)
	return q
}

func (t *InfraBaseTopology) SelectTopology(ctx context.Context) error {
	q := t.SelectInfraTopologyQuery()

	log.Debug().Interface("SelectTopologyQuery", q.LogHeader(Sn))
	err := t.SelectSingleChartsResources(ctx, q)
	return err
}
