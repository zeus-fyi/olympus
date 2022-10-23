package create

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type Chart struct {
	autogen_bases.ChartPackages
}

const Sn = "Chart"

func insertChart(c Chart) string {
	sqlInsertStatement := fmt.Sprintf(
		`WITH cte_insert_chart AS (
					 INSERT INTO chart_packages(chart_name, chart_version, chart_description)
					 VALUES ('%s', '%s', '%v')
					 RETURNING chart_package_id
 				 ) SELECT chart_package_id FROM cte_insert_chart`,
		c.ChartName, c.ChartVersion, c.ChartDescription)
	return sqlInsertStatement
}

func (c *Chart) InsertChart(ctx context.Context, q sql_query_templates.QueryParams, chart Chart) error {
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))
	query := insertChart(chart)
	err := apps.Pg.QueryRow(ctx, query).Scan(&c.ChartPackageID)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))

}
