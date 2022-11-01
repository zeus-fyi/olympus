package create_infra

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

const Sn = "Infrastructure"

func (i *InfraBaseTopology) SelectInfraTopologyQuery() {
	var insertTopQuery, insertTopOrgUserQuery, insertTopChartQuery, insertTopInfraQuery sql_query_templates.SubCTE

	// TODO, use db but too much of pita to update right now
	var ts chronos.Chronos
	chartPackageID := ts.UnixTimeStampNow()
	i.ChartPackageID = chartPackageID

	insertTopQuery.QueryName = "cte_insert_topology"
	insertTopQuery.RawQuery = fmt.Sprintf(`INSERT INTO topologies(name) VALUES ('%s') RETURNING topology_id`, i.Name)

	insertTopOrgUserQuery.QueryName = "cte_insert_ou_topology"
	insertTopOrgUserQuery.RawQuery = fmt.Sprintf(`
		  INSERT INTO org_users_topologies(topology_id, org_id, user_id)
		  VALUES ((SELECT topology_id FROM %s), %d, %d)`,
		insertTopQuery.QueryName, i.OrgID, i.UserID)

	insertTopChartQuery.QueryName = "cte_insert_chart_pkg"
	insertTopChartQuery.RawQuery = fmt.Sprintf(`
		  INSERT INTO chart_packages(chart_package_id, chart_name, chart_version, chart_description)
		  VALUES (%d, '%s', '%s', '%v')
	     RETURNING chart_package_id`,
		i.ChartPackageID, i.Chart.ChartName, i.ChartVersion, i.ChartDescription.String)

	insertTopInfraQuery.QueryName = "cte_insert_topology_infrastructure_components"
	insertTopInfraQuery.RawQuery = fmt.Sprintf(`
		  INSERT INTO topology_infrastructure_components(topology_id, chart_package_id)
		  VALUES ((SELECT topology_id FROM '%s'), (SELECT chart_package_id FROM '%s')) `,
		insertTopQuery.QueryName, insertTopChartQuery.QueryName)

	subCTEs := []sql_query_templates.SubCTE{insertTopQuery, insertTopOrgUserQuery, insertTopChartQuery, insertTopInfraQuery}
	i.CTE.AppendSubCtes(subCTEs)
	i.InsertPackagesCTE()
}

func (i *InfraBaseTopology) InsertInfraBase(ctx context.Context) error {
	var q sql_query_templates.QueryParams
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))
	i.InsertPackagesCTE()
	q.RawQuery = i.CTE.GenerateChainedCTE()
	r, err := apps.Pg.Exec(ctx, q.RawQuery)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("InsertTopologyClass: %s, Rows Affected: %d", q.LogHeader(Sn), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
