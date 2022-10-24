package containers

import (
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

// InsertPodContainerGroupSQL will use the next_id distributed ID generator and select the container id
// value for subsequent subcomponent relationships of its element, should greatly simplify the insert logic
func (p *PodContainersGroup) InsertPodContainerGroupSQL() string {
	// container

	podSpecChildClassTypeID := p.PodSpecTemplate.GetPodSpecChildClassTypeID()
	contSubCTE := sql_query_templates.NewSubInsertCTE("cte_insert_containers")
	contSubCTE.TableName = "containers"
	contSubCTE.Fields = []string{"container_id", "container_name", "container_image_id", "container_version_tag", "container_platform_os", "container_repository", "container_image_pull_policy"}

	// ports
	portsSubCTE := sql_query_templates.NewSubInsertCTE("cte_insert_container_ports")
	portsSubCTE.TableName = "container_ports"
	portsSubCTE.Fields = []string{"port_id", "port_name", "container_port", "host_port"}
	portsRelationshipsSubCTE := sql_query_templates.NewSubInsertCTE("cte_containers_ports_relationship")
	portsRelationshipsSubCTE.TableName = "containers_ports"
	portsRelationshipsSubCTE.Fields = []string{"chart_subcomponent_child_class_type_id", "container_id", "port_id"}

	// env vars
	envVarsSubCTE := sql_query_templates.NewSubInsertCTE("cte_container_environmental_vars")
	envVarsSubCTE.TableName = "container_environmental_vars"
	envVarsSubCTE.Fields = []string{"env_id", "name", "value"}
	envVarsRelationshipsSubCTE := sql_query_templates.NewSubInsertCTE("cte_container_environmental_vars_relationships")
	envVarsRelationshipsSubCTE.TableName = "containers_environmental_vars"
	envVarsRelationshipsSubCTE.Fields = []string{"chart_subcomponent_child_class_type_id", "container_id", "env_id"}

	// vms
	contVmsSubCTE := sql_query_templates.NewSubInsertCTE("cte_containers_volume_mounts")
	contVmsSubCTE.TableName = "container_volume_mounts"
	contVmsSubCTE.Fields = []string{"volume_mount_id", "volume_mount_path", "volume_name"}
	contVmsRelationshipsSubCTE := sql_query_templates.NewSubInsertCTE("cte_containers_volume_mounts_relationships")
	contVmsRelationshipsSubCTE.TableName = "containers_volume_mounts"
	contVmsRelationshipsSubCTE.Fields = []string{"chart_subcomponent_child_class_type_id", "container_id", "volume_mount_id"}

	// podSpec for containers
	podSpecSubCTE := sql_query_templates.NewSubInsertCTE("cte_insert_spec_pod_template_containers")
	podSpecSubCTE.TableName = "chart_subcomponent_spec_pod_template_containers"
	podSpecSubCTE.Fields = []string{"chart_subcomponent_child_class_type_id", "container_id", "is_init_container", "container_sort_order"}

	// volumes for pod spec
	podSpecVolumesSubCTE := sql_query_templates.NewSubInsertCTE("cte_pod_spec_volumes")
	podSpecVolumesSubCTE.TableName = "volumes"
	podSpecVolumesSubCTE.Fields = []string{"volume_id", "volume_name", "volume_key_values_jsonb"}
	podSpecVolumesRelationshipSubCTE := sql_query_templates.NewSubInsertCTE("cte_pod_spec_containers_volumes")
	podSpecVolumesRelationshipSubCTE.TableName = "containers_volumes"
	podSpecVolumesRelationshipSubCTE.Fields = []string{"chart_subcomponent_child_class_type_id", "volume_id"}

	p.insertVolumes(podSpecChildClassTypeID, &podSpecVolumesSubCTE, &podSpecVolumesRelationshipSubCTE)

	ts := chronos.Chronos{}
	// TODO for now will just generate ids here, something more complex can come later
	sortOrderIndex := 0
	for _, cont := range p.Containers {
		cont.ProcessAndSetAmbiguousContainerFieldStatusAndSubfieldIds()
		c := cont.Metadata
		// should continue appending values to header
		// container
		c.ContainerID = ts.UnixTimeStampNow()
		contSubCTE.AddValues(c.ContainerID, c.ContainerName, c.ContainerImageID, c.ContainerVersionTag, c.ContainerPlatformOs, c.ContainerRepository, c.ContainerImagePullPolicy)

		// pod spec to link container
		podSpecSubCTE.AddValues(podSpecChildClassTypeID, c.ContainerID, cont.IsInitContainer, sortOrderIndex)

		// ports
		p.getContainerPortsValuesForInsert(c.ContainerImageID, &portsSubCTE)
		p.getContainerPortsHeaderRelationshipValues(podSpecChildClassTypeID, c.ContainerImageID, &portsRelationshipsSubCTE)

		// env vars
		p.getInsertContainerEnvVarsValues(c.ContainerImageID, &envVarsSubCTE)
		p.getContainerEnvVarRelationshipValues(podSpecChildClassTypeID, c.ContainerImageID, &envVarsRelationshipsSubCTE)

		// vms
		p.insertContainerVolumeMountsValues(podSpecChildClassTypeID, c.ContainerImageID, &contVmsSubCTE, &contVmsRelationshipsSubCTE)
		sortOrderIndex += 1
	}

	cteExpr := sql_query_templates.CTE{
		Name: "InsertPodContainerGroupSQL",
		SubCTEs: []sql_query_templates.SubCTE{
			// container and podSpec template relationship
			contSubCTE,
			podSpecSubCTE,
			// ports
			portsSubCTE,
			portsRelationshipsSubCTE,
			// env vars
			envVarsSubCTE,
			envVarsRelationshipsSubCTE,
			// vms
			contVmsSubCTE,
			contVmsRelationshipsSubCTE,
			// vols
			podSpecVolumesSubCTE,
			podSpecVolumesRelationshipSubCTE,
		},
	}
	query := cteExpr.MultiLevelValuesCTEStringBuilderSQL()
	return query
}
