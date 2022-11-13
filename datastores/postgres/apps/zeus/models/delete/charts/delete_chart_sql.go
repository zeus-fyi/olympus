package delete_charts

import (
	"fmt"

	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

// DeleteChartQuery Not production ready
func DeleteChartQuery(q sql_query_templates.QueryParams) string {
	header := ""
	suffix := "SELECT true"
	switch q.QueryName {
	case "DeleteChart":
		header =
			`WITH cte_chart_packages AS (
				SELECT cpk.chart_package_id, cpk.chart_subcomponent_parent_class_type_id, cct.chart_subcomponent_child_class_type_id
				FROM chart_package_components cpk
				INNER JOIN chart_subcomponent_child_class_types cct ON cct.chart_subcomponent_parent_class_type_id = cpk.chart_subcomponent_parent_class_type_id
				WHERE chart_package_id = $1
			),`
	case "DeleteTopologyAndChart":
		// TODO, requires deleting lots of other things too..
	case "DeleteChartFromTopology":
		header = `
			WITH cte_chart_packages AS (
				SELECT cpk.chart_package_id AS chart_package_id, cpk.chart_subcomponent_parent_class_type_id, cct.chart_subcomponent_child_class_type_id
				FROM topology_infrastructure_components tic
				INNER JOIN topologies top ON top.topology_id = tic.topology_id
				INNER JOIN org_users_topologies out ON out.topology_id = tic.topology_id
				INNER JOIN chart_package_components cpk ON cpk.chart_package_id = tic.chart_package_id
				INNER JOIN chart_subcomponent_child_class_types cct ON cct.chart_subcomponent_parent_class_type_id = cpk.chart_subcomponent_parent_class_type_id
				WHERE tic.topology_id = $1 AND org_id = $2 AND user_id = $3
				LIMIT 1
		),`
	default:
		return ""
	}
	query := fmt.Sprintf(`
		%s
		cte_child_values AS (
			SELECT cv.chart_subcomponent_child_class_type_id
			FROM cte_chart_packages cpc 
			INNER JOIN chart_subcomponents_child_values cv ON cv.chart_subcomponent_child_class_type_id = cpc.chart_subcomponent_child_class_type_id
		), cte_ps_c AS (
			SELECT ps.container_id AS container_id, cpc.chart_subcomponent_child_class_type_id AS chart_subcomponent_child_class_type_id
			FROM cte_chart_packages cpc
			INNER JOIN chart_subcomponent_spec_pod_template_containers ps ON ps.chart_subcomponent_child_class_type_id = cpc.chart_subcomponent_child_class_type_id
		), cte_delete_volumes AS (
			SELECT cv.volume_id, ps.chart_subcomponent_child_class_type_id
			FROM cte_ps_c ps 
			INNER JOIN containers_volumes cv ON cv.chart_subcomponent_child_class_type_id = ps.chart_subcomponent_child_class_type_id
			INNER JOIN volumes v ON v.volume_id = cv.volume_id
		), cte_delete AS (
			SELECT psc.container_id, 
				cctx.container_security_context_id, 
				ccr.compute_resources_id, 
				cca.command_args_id,
				cprs.probe_id, 
				csports.port_id,
				csvms.volume_mount_id,
				csenvs.env_id
			FROM cte_ps_c psc
		  	LEFT JOIN containers cont ON cont.container_id = psc.	container_id
			LEFT JOIN containers_security_context cctx ON cctx.container_id = psc.container_id
			LEFT JOIN container_security_context csctx ON csctx.container_security_context_id = cctx.container_security_context_id
			LEFT JOIN containers_compute_resources cscr ON cscr.container_id = psc.container_id
			LEFT JOIN container_compute_resources AS ccr ON ccr.compute_resources_id = cscr.compute_resources_id 
			LEFT JOIN containers_command_args csca ON csca.container_id = psc.container_id
			LEFT JOIN container_command_args cca ON cca.command_args_id = csca.command_args_id
			LEFT JOIN containers_probes csprs ON csprs.container_id = psc.container_id
			LEFT JOIN container_probes cprs ON cprs.probe_id = csprs.probe_id
			LEFT JOIN containers_ports csports ON csports.container_id = psc.container_id
			LEFT JOIN container_ports cports ON cports.port_id = csports.port_id
			LEFT JOIN containers_volume_mounts csvms ON csvms.container_id = psc.container_id
			LEFT JOIN container_volume_mounts csvm ON csvm.volume_mount_id = csvms.volume_mount_id
			LEFT JOIN containers_environmental_vars csenvs ON csenvs.container_id = psc.container_id
			LEFT JOIN container_environmental_vars cenvs ON cenvs.env_id = csenvs.env_id
		), cte_delete_cs_ctx AS (
		  DELETE FROM containers_security_context 
			WHERE container_id IN (SELECT container_id FROM cte_delete)
		), cte_delete_c_ctx AS (
		  DELETE FROM container_security_context 
			WHERE container_security_context_id IN (SELECT container_security_context_id FROM cte_delete)
		), cte_delete_cscr AS (
		  DELETE FROM containers_compute_resources 
			WHERE container_id IN (SELECT container_id FROM cte_delete)
		), cte_delete_ccr AS (
		  DELETE FROM container_compute_resources 
			WHERE compute_resources_id IN (SELECT compute_resources_id FROM cte_delete)
		), cte_delete_csca AS (
		  DELETE FROM containers_command_args 
			WHERE container_id IN (SELECT container_id FROM cte_delete)
		), cte_delete_cca AS (
		  DELETE FROM container_command_args 
			WHERE command_args_id IN (SELECT command_args_id FROM cte_delete)
		), cte_delete_csprs AS (
		  DELETE FROM containers_probes 
			WHERE container_id IN (SELECT container_id FROM cte_delete)
		), cte_delete_cprs AS (
		  DELETE FROM container_probes 
			WHERE probe_id IN (SELECT probe_id FROM cte_delete)
		), cte_delete_csports AS (
		  DELETE FROM containers_ports 
			WHERE container_id IN (SELECT container_id FROM cte_delete)
		), cte_delete_cports AS (
		  DELETE FROM container_ports 
			WHERE port_id IN (SELECT port_id FROM cte_delete)
		), cte_delete_csvms AS (
		  DELETE FROM containers_volume_mounts 
			WHERE container_id IN (SELECT container_id FROM cte_delete)
		), cte_delete_csvm AS (
		  DELETE FROM container_volume_mounts 
			WHERE volume_mount_id IN (SELECT volume_mount_id FROM cte_delete)
		), cte_delete_cenvs AS (
		  DELETE FROM containers_environmental_vars 
			WHERE chart_subcomponent_child_class_type_id IN (SELECT chart_subcomponent_child_class_type_id FROM cte_chart_packages)
		), cte_delete_envs AS (
		  DELETE FROM container_environmental_vars 
			WHERE env_id IN (SELECT env_id FROM cte_delete)	
		), cte_delete_cvols AS (
		  DELETE FROM containers_volumes 
			WHERE chart_subcomponent_child_class_type_id IN (SELECT chart_subcomponent_child_class_type_id FROM cte_delete_volumes)
		), cte_delete_vols AS (
		  DELETE FROM volumes 
			WHERE volume_id IN (SELECT volume_id FROM cte_delete_volumes)
		), cte_delete_conts AS (
		  DELETE FROM containers 
			WHERE container_id IN (SELECT container_id FROM cte_ps_c)
		), cte_delete_ps_conts AS (
		  DELETE FROM chart_subcomponent_spec_pod_template_containers 
			WHERE container_id IN (SELECT container_id FROM cte_ps_c)
		), cte_delete_cvs AS (
		  DELETE FROM chart_subcomponents_child_values 
			WHERE chart_subcomponent_child_class_type_id IN (SELECT chart_subcomponent_child_class_type_id FROM cte_child_values)
		), cte_delete_cvst AS (
		  DELETE FROM chart_subcomponent_child_class_types
			WHERE chart_subcomponent_child_class_type_id IN (SELECT chart_subcomponent_child_class_type_id FROM cte_chart_packages)
		), cte_delete_parents_relation_to_chart AS (
		  DELETE FROM chart_package_components
			WHERE chart_subcomponent_parent_class_type_id IN (SELECT chart_subcomponent_parent_class_type_id FROM cte_chart_packages)
		), cte_delete_parents AS (
		  DELETE FROM chart_subcomponent_parent_class_types
			WHERE chart_subcomponent_parent_class_type_id IN (SELECT chart_subcomponent_parent_class_type_id FROM cte_chart_packages)
		) %s`, header, suffix)

	return query
}
