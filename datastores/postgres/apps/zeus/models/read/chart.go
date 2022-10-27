package read

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	deployments "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/deployments"
	read_deployments "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/deployments"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/read_containers"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type Chart struct {
	autogen_bases.ChartPackages
	autogen_bases.ChartComponentResources

	*deployments.Deployment
}

const ModelName = "Chart"

func fetchChartQuery(chartID int) string {
	query := fmt.Sprintf(`
	WITH cte_chart_packages AS (
			SELECT chart_package_id, chart_name, chart_version, chart_description
			FROM chart_packages
			WHERE chart_package_id = %d
	), cte_chart_package_components AS (							
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
			SELECT pcagg.chart_package_id, pcagg.chart_component_kind_name,
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
				cs.container_image_pull_policy
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
					json_build_object('name', 	cvm.volume_name, 'mountPath', cvm.volume_mount_path)) as container_vol_mounts
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
				ports.port_id_json_array AS port_id_json_array
			FROM cte_chart_subcomponent_spec_pod_template_containers ps
			LEFT JOIN cte_probes AS pr ON pr.container_id_probes = ps.container_id
			LEFT JOIN cte_container_environmental_vars AS cenv ON cenv.env_container_id = ps.container_id
			LEFT JOIN cte_container_volume_mounts AS cvm ON cvm.volume_mounts_container_id = ps.container_id
			LEFT JOIN cte_ports AS ports ON ports.container_id_ports = ps.container_id
	)
	SELECT  (SELECT chart_package_id FROM cte_chart_packages ) AS chart_package_id,
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
				COALESCE(cagg.port_id_json_array, '{}'::jsonb) AS ports,
				COALESCE(cagg.env_vars, '{}'::jsonb) AS env_vars,
				COALESCE(cagg.probes, '{}'::jsonb) AS probes,
				COALESCE(cagg.container_vol_mounts, '{}'::jsonb) AS container_vol_mounts,
				COALESCE(v.pod_spec_volumes, '{}'::jsonb) AS pod_spec_volumes
	FROM cte_chart_kind_agg_to_parent_children ckagg 
	LEFT JOIN cte_chart_subcomponent_spec_pod_template_containers AS ps ON ps.chart_package_id = ckagg.chart_package_id
	LEFT JOIN cte_pod_spec_volumes AS v ON v.child_class_type_id_pod_spec_volumes = ps.chart_subcomponent_child_class_type_id_ps
	LEFT JOIN cte_containers_agg AS cagg ON cagg.container_id = ps.container_id
`, chartID)
	return query
}

func (c *Chart) SelectSingleChartsResources(ctx context.Context, q sql_query_templates.QueryParams) error {
	log.Debug().Interface("SelectQuery", q.LogHeader(ModelName))
	q.RawQuery = fetchChartQuery(c.ChartPackageID)
	rows, err := apps.Pg.Query(ctx, q.RawQuery)
	if err != nil {
		log.Err(err).Msg(q.LogHeader(ModelName))
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var podSpecVolumesStr string
		container := read_containers.NewContainer()

		var ckagg string
		rowErr := rows.Scan(&c.ChartPackageID, &c.ChartName, &c.ChartVersion, &c.ChartDescription, &c.ChartComponentKindName,
			&ckagg,
			&container.Metadata.ContainerID, &container.Metadata.ContainerName, &container.Metadata.ContainerImageID,
			&container.Metadata.ContainerVersionTag, &container.Metadata.ContainerPlatformOs, &container.Metadata.ContainerRepository,
			&container.Metadata.ContainerImagePullPolicy,
			&container.DB.Ports,
			&container.DB.EnvVar, &container.DB.Probes, &container.DB.ContainerVolumes,
			&podSpecVolumesStr,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(ModelName))
			return rowErr
		}
		if c.ChartComponentKindName == "Deployment" {
			if c.Deployment == nil {

				deployment := deployments.NewDeployment()
				_ = read_deployments.ParseDeploymentParentChildAggValues(ckagg)
				c.Deployment = &deployment
				if len(podSpecVolumesStr) > 0 {
					vs, vserr := read_containers.ParsePodSpecDBVolumesString(podSpecVolumesStr)
					if vserr != nil {
						log.Err(vserr).Msg(q.LogHeader(ModelName))
						return vserr
					}
					c.Deployment.Spec.Template.SetK8sPodSpecVolumes(vs)
				}
			}
		}

		if container.Metadata.ContainerID != 0 {
			cerr := container.ParseFields()
			if cerr != nil {
				return cerr
			}
			c.Spec.Template.AddContainer(container.Container)
		}
	}
	return nil
}
