package containers

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

const SelectDeploymentResourceID = "(SELECT chart_component_resource_id FROM chart_component_resources WHERE chart_component_kind_name = 'Deployment' AND chart_component_api_version = 'apps/v1')"

// InsertPodTemplateSpecContainersCTE will use the next_id distributed ID generator and select the container id
// value for subsequent subcomponent relationships of its element, should greatly simplify the insert logic
func (p *PodTemplateSpec) InsertPodTemplateSpecContainersCTE(chart *charts.Chart) sql_query_templates.CTE {
	// container
	ts := chronos.Chronos{}
	if p.GetPodSpecParentClassTypeID() == 0 {
		p.SetPodSpecParentClassTypeID(ts.UnixTimeStampNow())
	}
	if p.GetPodSpecChildClassTypeID() == 0 {
		p.SetPodSpecChildClassTypeID(ts.UnixTimeStampNow())
	}
	contPodSpecParentClassCTE := sql_query_templates.NewSubInsertCTE("cte_podSpecParentClassTypeCTE")
	contPodSpecParentClassCTE.TableName = p.ChartSubcomponentParentClassTypes.GetTableName()
	contPodSpecParentClassCTE.Columns = []string{"chart_package_id", "chart_component_resource_id", "chart_subcomponent_parent_class_type_id", "chart_subcomponent_parent_class_type_name"}
	contPodSpecParentClassCTE.AddValues(chart.ChartPackageID, SelectDeploymentResourceID, p.GetPodSpecParentClassTypeID(), p.ChartSubcomponentParentClassTypes.ChartSubcomponentParentClassTypeName)

	cpkAddParentClassTypeSubCTEs := common.AddParentClassToChartPackage(chart, p.GetPodSpecParentClassTypeID())

	agCct := autogen_bases.ChartSubcomponentChildClassTypes{}
	contSubChildClassCTE := sql_query_templates.NewSubInsertCTE("cte_podSpecSubChildClassCTE")
	contSubChildClassCTE.TableName = agCct.GetTableName()
	contSubChildClassCTE.Columns = []string{"chart_subcomponent_parent_class_type_id", "chart_subcomponent_child_class_type_id", "chart_subcomponent_child_class_type_name"}
	contSubChildClassCTE.AddValues(p.GetPodSpecParentClassTypeID(), p.GetPodSpecChildClassTypeID(), "PodTemplateSpecChild")

	contSubCTE := sql_query_templates.NewSubInsertCTE("cte_insert_containers")
	contSubCTE.TableName = "containers"
	contSubCTE.Columns = []string{"container_id", "container_name", "container_image_id", "container_version_tag", "container_platform_os", "container_repository", "container_image_pull_policy"}

	// ports
	portsSubCTE := sql_query_templates.NewSubInsertCTE("cte_insert_container_ports")
	portsSubCTE.TableName = "container_ports"
	portsSubCTE.Columns = []string{"port_id", "port_name", "container_port", "host_port"}
	portsRelationshipsSubCTE := sql_query_templates.NewSubInsertCTE("cte_containers_ports_relationship")
	portsRelationshipsSubCTE.TableName = "containers_ports"
	portsRelationshipsSubCTE.Columns = []string{"chart_subcomponent_child_class_type_id", "container_id", "port_id"}

	// env vars
	envVarsSubCTE := sql_query_templates.NewSubInsertCTE("cte_container_environmental_vars")
	envVarsSubCTE.TableName = "container_environmental_vars"
	envVarsSubCTE.Columns = []string{"env_id", "name", "value"}
	envVarsRelationshipsSubCTE := sql_query_templates.NewSubInsertCTE("cte_container_environmental_vars_relationships")
	envVarsRelationshipsSubCTE.TableName = "containers_environmental_vars"
	envVarsRelationshipsSubCTE.Columns = []string{"chart_subcomponent_child_class_type_id", "container_id", "env_id"}

	// vms
	contVmsSubCTE := sql_query_templates.NewSubInsertCTE("cte_containers_volume_mounts")
	contVmsSubCTE.TableName = "container_volume_mounts"
	contVmsSubCTE.Columns = []string{"volume_mount_id", "volume_mount_path", "volume_name"}
	contVmsRelationshipsSubCTE := sql_query_templates.NewSubInsertCTE("cte_containers_volume_mounts_relationships")
	contVmsRelationshipsSubCTE.TableName = "containers_volume_mounts"
	contVmsRelationshipsSubCTE.Columns = []string{"chart_subcomponent_child_class_type_id", "container_id", "volume_mount_id"}

	// podSpec for containersMapByImageID
	podSpecSubCTE := sql_query_templates.NewSubInsertCTE("cte_insert_spec_pod_template_containers")
	podSpecSubCTE.TableName = "chart_subcomponent_spec_pod_template_containers"
	podSpecSubCTE.Columns = []string{"chart_subcomponent_child_class_type_id", "container_id", "is_init_container", "container_sort_order"}

	// probes
	contP := autogen_bases.ContainerProbes{}
	probesSubCTE := sql_query_templates.NewSubInsertCTE("cte_container_probes")
	probesSubCTE.TableName = "container_probes"
	probesSubCTE.Columns = contP.GetTableColumns()
	conPRelationship := autogen_bases.ContainersProbes{}
	probesRelationshipsSubCTE := sql_query_templates.NewSubInsertCTE("cte_containers_probes_relationship")
	probesRelationshipsSubCTE.TableName = "containers_probes"
	probesRelationshipsSubCTE.Columns = conPRelationship.GetTableColumns()

	podSpecVolumesSubCTE, podSpecVolumesRelationshipSubCTE := p.insertVolumes()

	sortOrderIndex := 0
	containersMapByImageID := p.NewPodContainersMapForDB()

	var probeCTEs sql_query_templates.SubCTEs
	var probeRelationshipCTEs sql_query_templates.SubCTEs
	for _, cont := range containersMapByImageID {
		c := cont.Metadata
		// should continue appending values to header
		// container

		// child class type to link to pod spec
		contSubCTE.AddValues(cont.GetContainerID(), c.ContainerName, c.ContainerImageID, c.ContainerVersionTag, c.ContainerPlatformOs, c.ContainerRepository, c.ContainerImagePullPolicy)

		// pod spec to link container
		podSpecSubCTE.AddValues(p.GetPodSpecChildClassTypeID(), cont.GetContainerID(), cont.IsInitContainer, sortOrderIndex)

		// ports
		p.getContainerPortsValuesForInsert(containersMapByImageID, c.ContainerImageID, &portsSubCTE)
		p.getContainerPortsHeaderRelationshipValues(containersMapByImageID, c.ContainerImageID, &portsRelationshipsSubCTE)

		// env vars
		p.getInsertContainerEnvVarsValues(containersMapByImageID, c.ContainerImageID, &envVarsSubCTE)
		p.getContainerEnvVarRelationshipValues(containersMapByImageID, c.ContainerImageID, &envVarsRelationshipsSubCTE)

		// vms
		p.insertContainerVolumeMountsValues(containersMapByImageID, c.ContainerImageID, &contVmsSubCTE, &contVmsRelationshipsSubCTE)

		// probes
		p1, p2 := common.CreateProbeValueSubCTEs(cont.GetContainerID(), cont.Probes)
		probeCTEs = sql_query_templates.AppendSubCteSlices(probeCTEs, p1)
		probeRelationshipCTEs = sql_query_templates.AppendSubCteSlices(probeRelationshipCTEs, p2)
		sortOrderIndex += 1
	}

	cteExpr := sql_query_templates.CTE{
		Name: "InsertPodTemplateSpecContainersCTE",
		SubCTEs: []sql_query_templates.SubCTE{
			// container and podSpec template relationship
			contSubCTE,
			contPodSpecParentClassCTE,
			cpkAddParentClassTypeSubCTEs,
			contSubChildClassCTE,
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
	cteExpr.AppendSubCtes(probeCTEs)
	cteExpr.AppendSubCtes(probeRelationshipCTEs)

	return cteExpr
}
