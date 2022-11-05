package read_charts

import (
	"fmt"

	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func FetchChartQuery(q sql_query_templates.QueryParams) string {
	header := ""
	switch q.QueryName {
	case "SelectSingleChartsResources":
		header =
			`WITH cte_chart_packages AS (
				SELECT chart_package_id, chart_name, chart_version, chart_description
				FROM chart_packages
				WHERE chart_package_id = $1
			)`
	case "SelectInfraTopologyQuery":
		header = `
			WITH cte_chart_packages AS (
				SELECT cp.chart_package_id AS chart_package_id, cp.chart_name, cp.chart_version, cp.chart_description
				FROM topology_infrastructure_components tic
				INNER JOIN topologies top ON top.topology_id = tic.topology_id
				INNER JOIN org_users_topologies out ON out.topology_id = tic.topology_id
				INNER JOIN chart_packages cp ON cp.chart_package_id = tic.chart_package_id
				WHERE tic.topology_id = $1 AND org_id = $2 AND user_id = $3
				LIMIT 1
		)`
	default:
		return ""
	}

	query := fmt.Sprintf(`
	%s, cte_chart_package_components AS (							
			SELECT  
				cp.chart_package_id,
				cp.chart_name,
				cp.chart_version,
				cp.chart_description,
				csp.chart_subcomponent_parent_class_type_id,
				csp.chart_subcomponent_parent_class_type_name,
				cpr.chart_component_resource_id, 
				cpr.chart_component_kind_name,
				cpr.chart_component_api_version
			FROM cte_chart_packages cp
			LEFT JOIN chart_package_components AS cpc ON cpc.chart_package_id = cp.chart_package_id
			LEFT JOIN chart_subcomponent_parent_class_types AS csp ON csp.chart_subcomponent_parent_class_type_id = cpc.chart_subcomponent_parent_class_type_id
			LEFT JOIN chart_component_resources AS cpr ON cpr.chart_component_resource_id = csp.chart_component_resource_id
	), cte_chart_package_components_values AS (
			SELECT 
				cpc.chart_package_id,
				cpc.chart_component_kind_name AS chart_component_kind_name,
				cpc.chart_subcomponent_parent_class_type_name,
				jsonb_object_agg(cpc.chart_subcomponent_parent_class_type_name,
							 jsonb_build_array(
								json_build_object('chart_subcomponent_parent_class_types', json_build_object('chart_subcomponent_parent_class_type_name',cpc.chart_subcomponent_parent_class_type_name, 'chart_subcomponent_parent_class_type_id', cpc.chart_subcomponent_parent_class_type_id )),
								json_build_object('chart_subcomponent_child_class_types', json_build_object('chart_subcomponent_child_class_type_id', cct.chart_subcomponent_child_class_type_id, 'chart_subcomponent_child_class_type_name', cct.chart_subcomponent_child_class_type_name)),
								json_build_object('chart_subcomponents_child_values', json_build_object('chart_subcomponent_key_name', cv.chart_subcomponent_key_name, 'chart_subcomponent_value', cv.chart_subcomponent_value))
							)
				) AS parent_child_values_obj_agg
			FROM cte_chart_package_components cpc
			LEFT JOIN chart_subcomponent_child_class_types AS cct ON cct.chart_subcomponent_parent_class_type_id = cpc.chart_subcomponent_parent_class_type_id
			LEFT JOIN chart_subcomponents_child_values AS cv ON cv.chart_subcomponent_child_class_type_id = cct.chart_subcomponent_child_class_type_id
			WHERE chart_subcomponent_key_name IS NOT NULL
			GROUP BY cpc.chart_package_id, cpc.chart_component_kind_name, cpc.chart_subcomponent_parent_class_type_id, cpc.chart_subcomponent_parent_class_type_name, cct.chart_subcomponent_child_class_type_id, cct.chart_subcomponent_child_class_type_name, cv.chart_subcomponent_value, cv.chart_subcomponent_key_name
	), cte_chart_package_components_values_parent_child_agg AS (
			SELECT cpcv.chart_package_id, cpcv.chart_component_kind_name, jsonb_agg(parent_child_values_obj_agg) AS parent_child_agg
			FROM cte_chart_package_components_values cpcv
			GROUP BY chart_package_id, chart_component_kind_name, chart_subcomponent_parent_class_type_name
	) , cte_chart_kind_agg_to_parent_children AS (
			SELECT  pcagg.chart_package_id, pcagg.chart_component_kind_name,
					json_build_object(pcagg.chart_component_kind_name, jsonb_agg(parent_child_agg)) as ckagg
			FROM cte_chart_package_components_values_parent_child_agg pcagg
			GROUP BY chart_package_id, chart_component_kind_name
	), cte_chart_subcomponent_spec_pod_template_containers AS (
			SELECT 
				cpc.chart_package_id AS chart_package_id,
				cpc.chart_component_kind_name,
				cpc.chart_subcomponent_parent_class_type_id AS chart_subcomponent_parent_class_type_id_ps,
				cpc.chart_subcomponent_parent_class_type_name AS chart_subcomponent_parent_class_type_name_ps,
				cct.chart_subcomponent_child_class_type_id AS chart_subcomponent_child_class_type_id_ps, 
				cct.chart_subcomponent_child_class_type_name AS chart_subcomponent_child_class_type_name_ps, 
				cs.container_id,
				cs.container_name,
				cs.container_image_id,
				cs.container_version_tag, 
				cs.container_platform_os,
				cs.container_repository,
				cs.container_image_pull_policy,
				cs.is_init_container as is_init_container
			FROM cte_chart_package_components cpc
			LEFT JOIN chart_subcomponent_child_class_types AS cct ON cct.chart_subcomponent_parent_class_type_id = cpc.chart_subcomponent_parent_class_type_id
			LEFT JOIN chart_subcomponent_spec_pod_template_containers AS ps ON ps.chart_subcomponent_child_class_type_id = cct.chart_subcomponent_child_class_type_id
			LEFT JOIN containers cs ON cs.container_id = ps.container_id
			LEFT JOIN chart_subcomponents_child_values AS cv ON cv.chart_subcomponent_child_class_type_id = cct.chart_subcomponent_child_class_type_id
			WHERE chart_subcomponent_key_name IS NULL
	), cte_container_environmental_vars AS (
			SELECT  cenv.container_id AS env_container_id, 
					jsonb_object_agg(cv.name, cv.value) AS env_vars
			FROM containers_environmental_vars cenv
			INNER JOIN cte_chart_subcomponent_spec_pod_template_containers AS pst ON pst.chart_subcomponent_child_class_type_id_ps = cenv.chart_subcomponent_child_class_type_id
			LEFT JOIN container_environmental_vars AS cv ON cv.env_id = cenv.env_id
			GROUP BY  cenv.container_id
	), cte_container_volume_mounts AS (
			SELECT  ps.container_id AS volume_mounts_container_id,
					jsonb_object_agg(cvm.volume_mount_id,
					json_build_object('name', 	cvm.volume_name, 'mountPath', cvm.volume_mount_path, 'subPath', cvm.volume_sub_path, 'readOnly', cvm.volume_read_only)) as container_vol_mounts
			FROM cte_chart_subcomponent_spec_pod_template_containers ps
			LEFT JOIN containers_volume_mounts AS csvm ON csvm.container_id = ps.container_id
			INNER JOIN container_volume_mounts AS cvm ON cvm.volume_mount_id = csvm.volume_mount_id
			GROUP BY ps.container_id
	), cte_pod_spec_volumes  AS (
			SELECT  csvm.chart_subcomponent_child_class_type_id AS child_class_type_id_pod_spec_volumes,
					jsonb_object_agg(volume_name, volume_key_values_jsonb) AS pod_spec_volumes
			FROM containers_volumes csvm
			LEFT JOIN cte_chart_subcomponent_spec_pod_template_containers AS cp ON cp.chart_subcomponent_child_class_type_id_ps = csvm.chart_subcomponent_child_class_type_id
			INNER JOIN volumes AS v ON v.volume_id = csvm.volume_id
			GROUP BY csvm.chart_subcomponent_child_class_type_id
	), cte_probes AS (
			SELECT					
				csp.container_id AS container_id_probes,
				jsonb_object_agg(probe_type, probe_key_values_jsonb) AS probes
			FROM containers_probes csp
			INNER JOIN cte_chart_subcomponent_spec_pod_template_containers AS cs ON cs.container_id = csp.container_id
			LEFT JOIN container_probes AS cpr ON cpr.probe_id = csp.probe_id 
			GROUP BY csp.container_id
	), cte_cmd_args AS (
			SELECT					
				csp.container_id AS cmd_container_id,
				jsonb_object_agg('cmdArgs',	json_build_object('command', command_values, 'args', args_values)) AS cmd_args
			FROM containers_command_args csp
			INNER JOIN cte_chart_subcomponent_spec_pod_template_containers AS cs ON cs.container_id = csp.container_id
			LEFT JOIN container_command_args AS cpr ON cpr.command_args_id = csp.command_args_id 
			GROUP BY csp.container_id
	), cte_comp_res AS (
		SELECT					
			csp.container_id AS comp_res_container_id,
			jsonb_object_agg('computeResources',
				json_build_object(
				'compute_resources_cpu_request', compute_resources_cpu_request, 
				'compute_resources_cpu_limit', compute_resources_cpu_limit,
				'compute_resources_ram_request', compute_resources_ram_request,
				'compute_resources_ram_limit', compute_resources_ram_limit,
				'compute_resources_ephemeral_storage_request', compute_resources_ephemeral_storage_request,
				'compute_resources_ephemeral_storage_limit', compute_resources_ephemeral_storage_limit
				)
			) AS comp_res
		FROM containers_compute_resources csp
		INNER JOIN cte_chart_subcomponent_spec_pod_template_containers AS cs ON cs.container_id = csp.container_id
		LEFT JOIN container_compute_resources AS cpr ON cpr.compute_resources_id = csp.compute_resources_id 
		GROUP BY csp.container_id
	), cte_ports AS (
			SELECT
				csp.container_id AS container_id_ports,
				jsonb_object_agg( 
				cp.port_id,
				json_build_object('name', cp.port_name, 'containerPort', cp.container_port, 'protocol', cp.port_protocol)) as port_id_json_array
			FROM containers_ports csp
			INNER JOIN cte_chart_subcomponent_spec_pod_template_containers AS cs ON cs.container_id = csp.container_id
			LEFT JOIN container_ports AS cp ON cp.port_id = csp.port_id 
			GROUP BY csp.container_id
	), cte_containers_agg AS (
			SELECT 
				ps.chart_subcomponent_parent_class_type_id_ps AS chart_subcomponent_parent_class_type_id_cagg,
				ps.chart_subcomponent_parent_class_type_name_ps AS chart_subcomponent_parent_class_type_name_cagg,
				ps.container_id, 
				env_vars,
				probes,
				container_vol_mounts AS container_vol_mounts,
				ports.port_id_json_array AS port_id_json_array,
				cmda.cmd_args AS cmd_args,
				compr.comp_res AS comp_res
			FROM cte_chart_subcomponent_spec_pod_template_containers ps
			LEFT JOIN cte_probes AS pr ON pr.container_id_probes = ps.container_id
			LEFT JOIN cte_container_environmental_vars AS cenv ON cenv.env_container_id = ps.container_id
			LEFT JOIN cte_container_volume_mounts AS cvm ON cvm.volume_mounts_container_id = ps.container_id
			LEFT JOIN cte_ports AS ports ON ports.container_id_ports = ps.container_id
			LEFT JOIN cte_cmd_args AS cmda ON cmda.cmd_container_id = ps.container_id
			LEFT JOIN cte_comp_res AS compr ON compr.comp_res_container_id = ps.container_id
	)
	SELECT  	(SELECT chart_package_id FROM cte_chart_packages ) AS chart_package_id,
				(SELECT chart_name FROM cte_chart_packages ) AS chart_name,
				(SELECT chart_version FROM cte_chart_packages ) AS chart_version,
				(SELECT chart_description FROM cte_chart_packages ) AS chart_description,
				ckagg.chart_component_kind_name,
				ckagg.ckagg AS ckagg,
				COALESCE(ps.container_id, 0) AS container_id,
				COALESCE(ps.container_name, '') AS container_name,
				COALESCE(ps.container_image_id,'') AS container_image_id,
				COALESCE(ps.container_version_tag, '') AS container_version_tag,
				COALESCE(ps.container_platform_os, '') AS container_platform_os,
				COALESCE(ps.container_repository,'') AS container_repository,
				COALESCE(ps.container_image_pull_policy,'') AS container_image_pull_policy,
				COALESCE(ps.is_init_container,'false') AS is_init_container,
				COALESCE(cagg.port_id_json_array, '{}'::jsonb) AS ports,
				COALESCE(cagg.env_vars, '{}'::jsonb) AS env_vars,
				COALESCE(cagg.probes, '{}'::jsonb) AS probes,
				COALESCE(cagg.container_vol_mounts, '{}'::jsonb) AS container_vol_mounts,
				COALESCE(cagg.cmd_args, '{}'::jsonb) AS cmd_args,
				COALESCE(cagg.comp_res, '{}'::jsonb) AS comp_res,
				COALESCE(v.pod_spec_volumes, '{}'::jsonb) AS pod_spec_volumes
	FROM cte_chart_kind_agg_to_parent_children ckagg 
	LEFT JOIN cte_chart_subcomponent_spec_pod_template_containers AS ps ON ps.chart_package_id = ckagg.chart_package_id
	LEFT JOIN cte_pod_spec_volumes AS v ON v.child_class_type_id_pod_spec_volumes = ps.chart_subcomponent_child_class_type_id_ps
	LEFT JOIN cte_containers_agg AS cagg ON cagg.container_id = ps.container_id
	ORDER BY container_id DESC
`, header)
	return query
	// TODO fix needing to order by container
}
