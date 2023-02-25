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

	var ts chronos.Chronos
	chartPackageID := ts.UnixTimeStampNow()
	i.Packages.ChartPackageID = chartPackageID
	i.ChartPackageID = chartPackageID

	insertTopQuery.QueryName = "cte_insert_topology"
	insertTopQuery.RawQuery = fmt.Sprintf(`INSERT INTO topologies(name) VALUES ($1) RETURNING topology_id`)
	i.CTE.ReturnSQLStatement = fmt.Sprintf("SELECT topology_id FROM %s", insertTopQuery.QueryName)

	i.CTE.Params = append(i.CTE.Params, i.Name)
	insertTopOrgUserQuery.QueryName = "cte_insert_ou_topology"
	insertTopOrgUserQuery.RawQuery = fmt.Sprintf(`
		  INSERT INTO org_users_topologies(topology_id, org_id, user_id)
		  VALUES ((SELECT topology_id FROM %s), $2, $3)`,
		insertTopQuery.QueryName)
	i.CTE.Params = append(i.CTE.Params, i.OrgID, i.UserID)

	insertTopChartQuery.QueryName = "cte_insert_chart_pkg"
	insertTopChartQuery.RawQuery = `
		  INSERT INTO chart_packages(chart_package_id, chart_name, chart_version, chart_description)
		  VALUES ($4, $5, $6, $7)
	     RETURNING chart_package_id`
	i.CTE.Params = append(i.CTE.Params, i.ChartPackageID, i.ChartName, i.ChartVersion, i.ChartDescription)

	insertTopInfraQuery.QueryName = "cte_insert_topology_infrastructure_components"

	if len(i.Tag) <= 0 {
		i.Tag = "latest"
	}

	if len(i.SkeletonBaseName) != 0 {
		if len(i.ClusterClassName) > 0 && len(i.ComponentBaseName) > 0 {
			insertTopInfraQuery.RawQuery = fmt.Sprintf(`
			  INSERT INTO topology_infrastructure_components(topology_id, chart_package_id, topology_skeleton_base_id, tag)
			  VALUES ((SELECT topology_id FROM %s), (SELECT chart_package_id FROM %s), (SELECT topology_skeleton_base_id
																						FROM topology_skeleton_base_components sb
																						INNER JOIN topology_base_components tb ON tb.topology_base_component_id = sb.topology_base_component_id
																						INNER JOIN topology_system_components sys ON sys.topology_system_component_id = tb.topology_system_component_id
																						WHERE sb.org_id = $8 AND sb.topology_skeleton_base_name = $9 AND sys.topology_system_component_name = $10 AND tb.topology_base_name = $11 ORDER BY topology_skeleton_base_id DESC LIMIT 1), $12) `,
				insertTopQuery.QueryName, insertTopChartQuery.QueryName)
			i.CTE.Params = append(i.CTE.Params, i.OrgID, i.SkeletonBaseName, i.ClusterClassName, i.ComponentBaseName, i.Tag)
		} else if len(i.ClusterClassName) > 0 {
			insertTopInfraQuery.RawQuery = fmt.Sprintf(`
			  INSERT INTO topology_infrastructure_components(topology_id, chart_package_id, topology_skeleton_base_id, tag)
			  VALUES ((SELECT topology_id FROM %s), (SELECT chart_package_id FROM %s), (SELECT topology_skeleton_base_id
																						FROM topology_skeleton_base_components sb
																						INNER JOIN topology_base_components tb ON tb.topology_base_component_id = sb.topology_base_component_id
																						INNER JOIN topology_system_components sys ON sys.topology_system_component_id = tb.topology_system_component_id
																						WHERE sb.org_id = $8 AND sb.topology_skeleton_base_name = $9 AND sys.topology_system_component_name = $10 ORDER BY topology_skeleton_base_id DESC LIMIT 1), $11) `,
				insertTopQuery.QueryName, insertTopChartQuery.QueryName)
			i.CTE.Params = append(i.CTE.Params, i.OrgID, i.SkeletonBaseName, i.ClusterClassName, i.Tag)
		} else {
			insertTopInfraQuery.RawQuery = fmt.Sprintf(`
		  INSERT INTO topology_infrastructure_components(topology_id, chart_package_id, topology_skeleton_base_id, tag)
		  VALUES ((SELECT topology_id FROM %s), (SELECT chart_package_id FROM %s), (SELECT topology_skeleton_base_id FROM topology_skeleton_base_components WHERE org_id = $8 AND topology_skeleton_base_name = $9 ORDER BY topology_skeleton_base_id DESC LIMIT 1), $10) `,
				insertTopQuery.QueryName, insertTopChartQuery.QueryName)
			i.CTE.Params = append(i.CTE.Params, i.OrgID, i.SkeletonBaseName, i.Tag)
		}
	} else {
		insertTopInfraQuery.RawQuery = fmt.Sprintf(`
		  INSERT INTO topology_infrastructure_components(topology_id, chart_package_id, topology_skeleton_base_id, tag)
		  VALUES ((SELECT topology_id FROM %s), (SELECT chart_package_id FROM %s), $8, $9)`,
			insertTopQuery.QueryName, insertTopChartQuery.QueryName)
		i.CTE.Params = append(i.CTE.Params, ts.UnixTimeStampNow(), i.Tag)
	}

	subCTEs := []sql_query_templates.SubCTE{insertTopQuery, insertTopOrgUserQuery, insertTopChartQuery, insertTopInfraQuery}
	i.CTE.AppendSubCtes(subCTEs)
	i.InsertPackagesCTE()
}

func (i *InfraBaseTopology) InsertInfraBase(ctx context.Context) error {
	var q sql_query_templates.QueryParams
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))
	i.SelectInfraTopologyQuery()
	q.RawQuery = i.CTE.GenerateChainedCTE()
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, i.CTE.Params...).Scan(&i.TopologyID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
