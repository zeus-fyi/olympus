package read

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/containers"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/deployments"
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
			LEFT JOIN chart_subcomponent_parent_class_types AS csp ON csp.chart_subcomponent_parent_class_type_id =  cpc.chart_subcomponent_parent_class_type_id
			LEFT JOIN chart_component_resources AS cpr ON cpr.chart_component_resource_id = csp.chart_component_resource_id
	), cte_chart_package_components_values AS (
				SELECT 
				cpc.chart_component_kind_name AS chart_component_kind_name,
				cpc.chart_subcomponent_parent_class_type_id AS chart_subcomponent_parent_class_type_id, 
				cpc.chart_subcomponent_parent_class_type_name AS chart_subcomponent_parent_class_type_name,
				cct.chart_subcomponent_child_class_type_id AS chart_subcomponent_child_class_type_id,
				cct.chart_subcomponent_child_class_type_name AS chart_subcomponent_child_class_type_name, 
				jsonb_object_agg(cv.chart_subcomponent_key_name, cv.chart_subcomponent_value) AS child_values
			FROM cte_chart_package_components cpc
			LEFT JOIN chart_subcomponent_child_class_types AS cct ON cct.chart_subcomponent_parent_class_type_id = cpc.chart_subcomponent_parent_class_type_id
			LEFT JOIN chart_subcomponents_child_values AS cv ON cv.chart_subcomponent_child_class_type_id = cct.chart_subcomponent_child_class_type_id
			WHERE chart_subcomponent_key_name IS NOT NULL
			GROUP BY cpc.chart_component_kind_name, cpc.chart_subcomponent_parent_class_type_id, cpc.chart_subcomponent_parent_class_type_name, cct.chart_subcomponent_child_class_type_id, cct.chart_subcomponent_child_class_type_name
	 ), cte_chart_subcomponent_spec_pod_template_containers AS (
			SELECT 
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
	), cte_container_environmental_vars AS (
			SELECT  cenv.container_id AS env_container_id, 
					jsonb_object_agg(cv.name, cv.value) AS env_vars
			FROM containers_environmental_vars cenv
			INNER JOIN cte_chart_subcomponent_spec_pod_template_containers AS pst ON pst.chart_subcomponent_child_class_type_id_ps = cenv.chart_subcomponent_child_class_type_id
			LEFT JOIN container_environmental_vars AS cv ON cv.env_id = cenv.env_id
			GROUP BY  cenv.container_id
	), cte_container_volume_mounts AS (
			SELECT 	ps.container_id AS volume_mounts_container_id,
					cvm.volume_name, 
					cvm.volume_mount_path
			FROM cte_chart_subcomponent_spec_pod_template_containers ps
			LEFT JOIN containers_volume_mounts AS csvm ON csvm.container_id = ps.container_id
			LEFT JOIN container_volume_mounts AS cvm ON cvm.volume_mount_id = csvm.volume_mount_id
	), cte_pod_spec_volumes  AS (
			SELECT  csvm.chart_subcomponent_child_class_type_id AS child_class_type_id_pod_spec_volumes,
					jsonb_object_agg(volume_name, volume_key_values_jsonb) AS pod_spec_volumes
			FROM containers_volumes csvm
			INNER JOIN cte_chart_subcomponent_spec_pod_template_containers AS cp ON cp.chart_subcomponent_child_class_type_id_ps = csvm.chart_subcomponent_child_class_type_id
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
				cp.port_name AS port_name,
				cp.container_port AS container_port,
				cp.port_protocol AS port_protcol
			FROM containers_ports csp
			INNER JOIN cte_chart_subcomponent_spec_pod_template_containers AS cs ON cs.container_id = csp.container_id
			LEFT JOIN container_ports AS cp ON cp.port_id = csp.port_id 
	), cte_containers_agg AS (
			SELECT 
				ps.chart_subcomponent_parent_class_type_id_ps AS chart_subcomponent_parent_class_type_id_cagg,
				ps.chart_subcomponent_parent_class_type_name_ps AS chart_subcomponent_parent_class_type_name_cagg,
				ps.container_id, 
				env_vars,
				probes,
				volume_name, 
				volume_mount_path,
				ports.port_name AS port_name,
				ports.container_port AS container_port,
				ports.port_protcol AS port_protcol
			FROM cte_chart_subcomponent_spec_pod_template_containers ps
			LEFT JOIN cte_probes AS pr ON pr.container_id_probes = ps.container_id
			LEFT JOIN cte_container_environmental_vars AS cenv ON cenv.env_container_id = ps.container_id
			LEFT JOIN cte_container_volume_mounts AS cvm ON cvm.volume_mounts_container_id = ps.container_id
			LEFT JOIN cte_ports AS ports ON ports.container_id_ports = ps.container_id
	), cte_chart_package_components_values_agg AS (
			SELECT 
				chart_subcomponent_parent_class_type_id AS chart_subcomponent_parent_class_type_id_agg, 
				chart_subcomponent_parent_class_type_name AS chart_subcomponent_parent_class_type_name_agg,
				chart_subcomponent_child_class_type_id AS chart_subcomponent_child_class_type_id_agg,
				chart_subcomponent_child_class_type_name AS chart_subcomponent_child_class_type_name_agg,
				jsonb_object_agg(chart_subcomponent_parent_class_type_name,  child_values) AS child_values
			FROM cte_chart_package_components_values
			GROUP BY  chart_subcomponent_parent_class_type_id, chart_subcomponent_child_class_type_id, chart_subcomponent_parent_class_type_name, chart_subcomponent_child_class_type_name
	)  	
	SELECT  (SELECT chart_package_id FROM cte_chart_packages ) AS chart_package_id,
			(SELECT chart_name FROM cte_chart_packages ) AS chart_name,
			(SELECT chart_version FROM cte_chart_packages ) AS chart_version,
			(SELECT chart_description FROM cte_chart_packages ) AS chart_description,
			chart_component_kind_name,
			COALESCE(chart_subcomponent_parent_class_type_id_agg, chart_subcomponent_parent_class_type_id_ps) AS chart_subcomponent_parent_class_type_id,
			COALESCE(chart_subcomponent_parent_class_type_name_agg, chart_subcomponent_parent_class_type_name_ps) AS chart_subcomponent_parent_class_type_name,
			COALESCE(chart_subcomponent_child_class_type_id_agg, chart_subcomponent_child_class_type_id_ps) AS chart_subcomponent_child_class_type_id,
			COALESCE(chart_subcomponent_child_class_type_name_agg, chart_subcomponent_child_class_type_name_ps) AS chart_subcomponent_child_class_type_name,
			COALESCE(ps.container_id, 0) AS container_id,
			COALESCE(ps.container_name, '') AS container_name,
			COALESCE(ps.container_image_id,'') AS container_image_id,
			COALESCE(ps.container_version_tag, '') AS container_version_tag,
			COALESCE(ps.container_platform_os, '') AS container_platform_os,
			COALESCE(ps.container_repository,'') AS container_repository,
			COALESCE(ps.container_image_pull_policy,'') AS container_image_pull_policy,
			COALESCE(cagg.port_name, '{}') AS port_name,
			COALESCE(cagg.container_port, 0) AS container_port,
			COALESCE(cagg.port_protcol, '{}') AS port_protcol,
			COALESCE(cagg.env_vars, '{}'::jsonb) AS env_vars,
			COALESCE(cagg.probes, '{}'::jsonb) AS probes,
			COALESCE(cagg.volume_name, '') AS volume_name,
			COALESCE(cagg.volume_mount_path, '') AS volume_mount_path,
			COALESCE(v.pod_spec_volumes, '{}'::jsonb) AS pod_spec_volumes
	FROM cte_chart_subcomponent_spec_pod_template_containers ps
	LEFT JOIN cte_containers_agg AS cagg ON cagg.container_id = ps.container_id
	LEFT JOIN cte_chart_package_components_values_agg AS agg ON agg.chart_subcomponent_child_class_type_id_agg = ps.chart_subcomponent_child_class_type_id_ps
	LEFT JOIN cte_pod_spec_volumes AS v ON v.child_class_type_id_pod_spec_volumes = ps.chart_subcomponent_child_class_type_id_ps
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
	//var podTemplateSpec containers.PodTemplateSpec
	parentClassSlice := autogen_bases.ChartSubcomponentParentClassTypesSlice{}
	childClassSlice := autogen_bases.ChartSubcomponentChildClassTypesSlice{}

	for rows.Next() {
		parentClassElement := autogen_bases.ChartSubcomponentParentClassTypes{}
		childClassElement := autogen_bases.ChartSubcomponentChildClassTypes{}
		volumePtr := autogen_bases.Volumes{}
		container := containers.NewContainer()

		rowErr := rows.Scan(&c.ChartPackageID, &c.ChartName, &c.ChartVersion, &c.ChartDescription, &c.ChartComponentKindName,
			&parentClassElement.ChartSubcomponentParentClassTypeID, &parentClassElement.ChartSubcomponentParentClassTypeName,
			&childClassElement.ChartSubcomponentChildClassTypeID, &childClassElement.ChartSubcomponentChildClassTypeName,

			&container.Metadata.ContainerID, &container.Metadata.ContainerName, &container.Metadata.ContainerImageID,
			&container.Metadata.ContainerVersionTag, &container.Metadata.ContainerPlatformOs, &container.Metadata.ContainerRepository,
			&container.Metadata.ContainerImagePullPolicy,
			&container.DB.PortName, &container.DB.PortNumber, &container.DB.PortProtocol,
			&container.DB.EnvVar, &container.DB.Probes, &container.DB.VolumeName, &container.DB.VolumePath,
			&volumePtr.VolumeKeyValuesJSONb,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(ModelName))
			return rowErr
		}
		if c.ChartComponentKindName == "Deployment" {
			if c.Deployment == nil {
				deployment := deployments.Deployment{}
				c.Deployment = &deployment
				//if volumePtr.VolumeKeyValuesJSONb != "" {
				//	c.Deployment.Spec.Template.AddVolume()
				//}
			}

			// TODO add metadata, spec to deployment
			// TODO add container conversion json -> probes, env_vars
			// TODO add volume -> spec
			if container.Metadata.ContainerID != 0 {
				cerr := container.ParseFields()
				if cerr != nil {
					return cerr
				}
				c.Deployment.Spec.Template.AddContainer(container.Container)
			}
		}

		parentClassSlice = append(parentClassSlice, parentClassElement)
		childClassSlice = append(childClassSlice, childClassElement)
	}
	return nil
}
