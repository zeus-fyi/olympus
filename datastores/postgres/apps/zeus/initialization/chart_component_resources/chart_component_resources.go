package chart_component_resources

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type ChartComponentResources struct {
	autogen_bases.ChartComponentResources
}

const Sn = "ChartComponentResources"

func (c *ChartComponentResources) insertChartResource() string {
	sqlInsertStatement := fmt.Sprintf(
		`INSERT INTO chart_component_resources(chart_component_resource_id, chart_component_kind_name, chart_component_api_version)
 				 VALUES ('%d', '%s', '%s')`,
		c.ChartComponentResourceID, c.ChartComponentKindName, c.ChartComponentApiVersion)
	return sqlInsertStatement
}

func (c *ChartComponentResources) InsertChartResource(ctx context.Context, q sql_query_templates.QueryParams) error {
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))
	query := c.insertChartResource()
	_, err := apps.Pg.Exec(ctx, query)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
