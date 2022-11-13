package delete_chart_from_topology

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	delete_charts "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/delete/charts"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type DeleteChartFromTopology struct {
	TopologyID int
	OrgID      int
	UserID     int
	Deleted    bool
}

const Sn = "DeleteChartFromTopology"

func (c *DeleteChartFromTopology) getDeleteChartFromTopologyQuery() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "DeleteChartFromTopology"
	q.RawQuery = delete_charts.DeleteChartQuery(q)
	return q
}

func (c *DeleteChartFromTopology) DeleteChartFromTopology(ctx context.Context) error {
	q := c.getDeleteChartFromTopologyQuery()
	log.Debug().Interface("DeleteQuery:", q.LogHeader(Sn))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, c.TopologyID, c.OrgID, c.UserID).Scan(&c.Deleted)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
