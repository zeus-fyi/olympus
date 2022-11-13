package delete_charts

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type DeleteChart struct {
	ChartID int
	Deleted bool
}

const Sn = "DeleteChart"

func (c *DeleteChart) getDeleteChartQuery() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "DeleteChart"
	q.RawQuery = DeleteChartQuery(q)
	return q
}

func (c *DeleteChart) DeleteChart(ctx context.Context) error {
	q := c.getDeleteChartQuery()
	log.Debug().Interface("DeleteQuery:", q.LogHeader(Sn))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, c.ChartID).Scan(&c.Deleted)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
