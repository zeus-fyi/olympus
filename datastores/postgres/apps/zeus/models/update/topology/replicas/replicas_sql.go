package update_replicas

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/deployments"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/statefulsets"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func UpdateReplicaCountSQL(q sql_query_templates.QueryParams, replicaCount string) sql_query_templates.QueryParams {
	header := ""
	childUpdate := ""
	newVersion := fmt.Sprintf("version-%d", ts.UnixTimeStampNow())
	versionUpdate := ""
	switch q.QueryName {
	case "UpdateReplicaCountDeployment":
		header = getChartResourceTypes(deployments.DeploymentChartComponentResourceID)
		childUpdate = getUpdate(&q, "Spec", "replicas", replicaCount)
		versionUpdate = getUpdate(&q, "DeploymentParentMetadata", "version", newVersion)
	case "UpdateReplicaCountStatefulSet":
		header = getChartResourceTypes(statefulsets.StsChartComponentResourceID)
		childUpdate = getUpdate(&q, "Spec", "replicas", replicaCount)
		versionUpdate = getUpdate(&q, "StatefulSetParentMetadata", "version", newVersion)
	default:
		return q
	}
	q.RawQuery = fmt.Sprintf("%s %s %s SELECT true", header, childUpdate, versionUpdate)
	return q
}

func getUpdate(q *sql_query_templates.QueryParams, parentClassTypeName, key, value string) string {

	q.CTEQuery.Params = append(q.CTEQuery.Params, parentClassTypeName)
	c0 := len(q.CTEQuery.Params)
	q.CTEQuery.Params = append(q.CTEQuery.Params, key)
	c1 := len(q.CTEQuery.Params)
	q.CTEQuery.Params = append(q.CTEQuery.Params, value)
	c2 := len(q.CTEQuery.Params)
	q.CTEQuery.Params = append(q.CTEQuery.Params, key)
	c4 := len(q.CTEQuery.Params)
	getCteIDName := fmt.Sprintf("cte_get_field_to_update_%d", ts.UnixTimeStampNow())
	query := fmt.Sprintf(`, %s AS (
				SELECT cv.chart_subcomponent_child_class_type_id
				FROM chart_subcomponents_child_values cv
				INNER JOIN cte_chart_packages ccp ON ccp.chart_subcomponent_child_class_type_id = cv.chart_subcomponent_child_class_type_id
				WHERE ccp.chart_subcomponent_parent_class_type_name = $%d AND cv.chart_subcomponent_key_name = $%d
			), cte_update_%d AS (
				UPDATE chart_subcomponents_child_values
			  	SET chart_subcomponent_value = $%d
				WHERE chart_subcomponent_child_class_type_id IN (SELECT chart_subcomponent_child_class_type_id FROM %s) AND chart_subcomponent_key_name = $%d
			)`, getCteIDName, c0, c1, ts.UnixTimeStampNow(), c2, getCteIDName, c4)
	return query
}

func getChartResourceTypes(chartComponentResourceID int) string {
	header := fmt.Sprintf(`
			WITH cte_chart_packages AS (
				SELECT cpk.chart_package_id AS chart_package_id,
					   cpk.chart_subcomponent_parent_class_type_id,
					   cct.chart_subcomponent_child_class_type_id,
					   cct.chart_subcomponent_child_class_type_name,
 					   cct.chart_subcomponent_parent_class_type_id,
					   pct.chart_subcomponent_parent_class_type_name
				FROM topology_infrastructure_components tic
				INNER JOIN topologies top ON top.topology_id = tic.topology_id
				INNER JOIN org_users_topologies out ON out.topology_id = tic.topology_id
				INNER JOIN chart_package_components cpk ON cpk.chart_package_id = tic.chart_package_id
				INNER JOIN chart_subcomponent_child_class_types cct ON cct.chart_subcomponent_parent_class_type_id = cpk.chart_subcomponent_parent_class_type_id
				INNER JOIN chart_subcomponent_parent_class_types pct ON pct.chart_subcomponent_parent_class_type_id = cpk.chart_subcomponent_parent_class_type_id
				INNER JOIN chart_component_resources ccres ON ccres.chart_component_resource_id = pct.chart_component_resource_id
				WHERE tic.topology_id = $1 AND org_id = $2 AND user_id = $3 AND ccres.chart_component_resource_id = %d
			)`, chartComponentResourceID)

	return header
}
