package read

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type Chart struct {
	autogen_bases.ChartPackages
}

const ModelName = "Chart"

var query = `
	WITH cte_chart_packages AS (
			SELECT chart_package_id, chart_name, chart_version, chart_description
			FROM chart_packages
			ORDER BY chart_package_id DESC 
			LIMIT 1
	), cte_chart_package_components AS (
			SELECT  csp.chart_subcomponent_parent_class_type_id, csp.chart_subcomponent_parent_class_type_name
			FROM cte_chart_packages cp
			LEFT JOIN chart_package_components AS cpc ON cpc.chart_package_id = cp.chart_package_id
			LEFT JOIN chart_subcomponent_parent_class_types AS csp ON csp.chart_subcomponent_parent_class_type_id =  cpc.chart_subcomponent_parent_class_type_id
	), cte_chart_package_components_values AS (
			SELECT 
				   cct.chart_subcomponent_child_class_type_id, cct.chart_subcomponent_child_class_type_name, 
				   cv.chart_subcomponent_key_name, cv.chart_subcomponent_value
			FROM cte_chart_package_components cpc
			LEFT JOIN chart_subcomponent_child_class_types AS cct ON cct.chart_subcomponent_parent_class_type_id = cpc.chart_subcomponent_parent_class_type_id
			LEFT JOIN chart_subcomponents_child_values AS cv ON cv.chart_subcomponent_child_class_type_id = cct.chart_subcomponent_child_class_type_id
			LEFT JOIN chart_subcomponents_jsonb_child_values AS cvj ON cvj.chart_subcomponent_child_class_type_id = cct.chart_subcomponent_child_class_type_id
			LEFT JOIN chart_subcomponents_jsonb_child_values AS cvps ON cvps.chart_subcomponent_child_class_type_id = cct.chart_subcomponent_child_class_type_id
	),  cte_chart_subcomponent_spec_pod_template_containers AS (
			SELECT cs.container_id, cs.container_name, cs.container_image_id, cs.container_version_tag, cs.container_platform_os, cs.container_repository, cs.container_image_pull_policy
			FROM chart_subcomponent_spec_pod_template_containers ps
			INNER JOIN cte_chart_package_components_values AS cp ON cp.chart_subcomponent_child_class_type_id = ps.chart_subcomponent_child_class_type_id
			LEFT JOIN containers cs ON cs.container_id = ps.container_id
	), cte_container_environmental_vars AS (
			SELECT cenv.env_id, cv.name, cv.value
			FROM containers_environmental_vars cenv
			LEFT JOIN container_environmental_vars AS cv ON cv.env_id = cenv.env_id
	), cte_container_volume_mounts AS (
			SELECT ps.container_id, cvm.volume_name, cvm.volume_mount_path
			FROM containers_volume_mounts csvm
			LEFT JOIN cte_chart_subcomponent_spec_pod_template_containers AS ps ON ps.container_id = csvm.container_id
			LEFT JOIN container_volume_mounts AS cvm ON cvm.volume_mount_id = csvm.volume_mount_id
	), cte_containers_volumes  AS (
			SELECT  csvm.chart_subcomponent_child_class_type_id, jsonb_object_agg(volume_name, volume_key_values_jsonb)
			FROM containers_volumes csvm
			INNER JOIN cte_chart_package_components_values AS cp ON cp.chart_subcomponent_child_class_type_id = csvm.chart_subcomponent_child_class_type_id
			INNER JOIN volumes AS v ON v.volume_id = csvm.volume_id
			GROUP BY csvm.chart_subcomponent_child_class_type_id
	), cte_probes AS (
			SELECT					
			csp.container_id, jsonb_object_agg(probe_type, probe_key_values_jsonb)
			FROM containers_probes csp
			LEFT JOIN cte_chart_subcomponent_spec_pod_template_containers AS cs ON cs.container_id = csp.container_id
			LEFT JOIN container_probes AS cpr ON cpr.probe_id = csp.probe_id 
			GROUP BY csp.container_id
	)
	SELECT * FROM cte_chart_package_components
`

func (s *Chart) SelectChartResources(ctx context.Context, q sql_query_templates.QueryParams) error {
	log.Debug().Interface("SelectQuery", q.LogHeader(ModelName))
	rows, err := apps.Pg.Query(ctx, q.SelectQuery())
	if err != nil {
		log.Err(err).Msg(q.LogHeader(ModelName))
		return err
	}
	defer rows.Close()
	//var podTemplateSpec containers.PodTemplateSpec
	for rows.Next() {
		rowErr := rows.Scan()
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(ModelName))
			return rowErr
		}
		//selectedStructNameExamples = append(selectedStructNameExamples, se)
	}
	return nil
}
