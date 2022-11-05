package containers

import (
	"fmt"

	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

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

	contPodSpecParentClassCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_podSpecParentClassTypeCTE_%d", ts.UnixTimeStampNow()))
	contPodSpecParentClassCTE.TableName = p.ChartSubcomponentParentClassTypes.GetTableName()
	contPodSpecParentClassCTE.Columns = []string{"chart_package_id", "chart_component_resource_id", "chart_subcomponent_parent_class_type_id", "chart_subcomponent_parent_class_type_name"}
	contPodSpecParentClassCTE.AddValues(chart.ChartPackageID, p.ChartComponentResourceID, p.GetPodSpecParentClassTypeID(), p.ChartSubcomponentParentClassTypes.ChartSubcomponentParentClassTypeName)

	cpkAddParentClassTypeSubCTEs := common.AddParentClassToChartPackage(chart, p.GetPodSpecParentClassTypeID())

	templateMetadataCTE := common.CreateParentMetadataSubCTEs(chart, p.Metadata)
	agCct := autogen_bases.ChartSubcomponentChildClassTypes{}
	contSubChildClassCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_podSpecSubChildClassCTE_%d", ts.UnixTimeStampNow()))
	contSubChildClassCTE.TableName = agCct.GetTableName()
	contSubChildClassCTE.Columns = []string{"chart_subcomponent_parent_class_type_id", "chart_subcomponent_child_class_type_id", "chart_subcomponent_child_class_type_name"}
	contSubChildClassCTE.AddValues(p.GetPodSpecParentClassTypeID(), p.GetPodSpecChildClassTypeID(), "PodTemplateSpecChild")

	contSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_insert_containers_%d", ts.UnixTimeStampNow()))
	contSubCTE.TableName = "containers"
	contSubCTE.Columns = []string{"container_id", "container_name", "container_image_id", "container_version_tag", "container_platform_os", "container_repository", "container_image_pull_policy", "is_init_container"}

	// ports
	portsSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_insert_container_ports_%d", ts.UnixTimeStampNow()))
	portsSubCTE.TableName = "container_ports"
	portsSubCTE.Columns = []string{"port_id", "port_name", "container_port", "host_port"}
	portsRelationshipsSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_containers_ports_relationship_%d", ts.UnixTimeStampNow()))
	portsRelationshipsSubCTE.TableName = "containers_ports"
	portsRelationshipsSubCTE.Columns = []string{"chart_subcomponent_child_class_type_id", "container_id", "port_id"}

	// env vars
	envVarsSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_container_environmental_vars_%d", ts.UnixTimeStampNow()))
	envVarsSubCTE.TableName = "container_environmental_vars"
	envVarsSubCTE.Columns = []string{"env_id", "name", "value"}
	envVarsRelationshipsSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_container_environmental_vars_relationships_%d", ts.UnixTimeStampNow()))
	envVarsRelationshipsSubCTE.TableName = "containers_environmental_vars"
	envVarsRelationshipsSubCTE.Columns = []string{"chart_subcomponent_child_class_type_id", "container_id", "env_id"}

	// cmd args
	cmdArgsSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_container_command_args_%d", ts.UnixTimeStampNow()))
	cmdArgsSubCTE.TableName = "container_command_args"
	cmdArgsSubCTE.Columns = []string{"command_args_id", "command_values", "args_values"}
	cmdArgsRelationshipsSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_containers_command_args_relationships_%d", ts.UnixTimeStampNow()))
	cmdArgsRelationshipsSubCTE.TableName = "containers_command_args"
	cmdArgsRelationshipsSubCTE.Columns = []string{"command_args_id", "container_id"}

	// computeResources
	computeResourcesSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_container_compute_resources_%d", ts.UnixTimeStampNow()))
	computeResourcesSubCTE.TableName = "container_compute_resources"
	computeResourcesSubCTE.Columns = []string{"compute_resources_id", "compute_resources_cpu_request", "compute_resources_cpu_limit",
		"compute_resources_ram_request", "compute_resources_ram_limit", "compute_resources_ephemeral_storage_request", "compute_resources_ephemeral_storage_limit"}
	computeResourcesRelationshipsSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_containers_compute_resources_relationships_%d", ts.UnixTimeStampNow()))
	computeResourcesRelationshipsSubCTE.TableName = "containers_compute_resources"
	computeResourcesRelationshipsSubCTE.Columns = []string{"compute_resources_id", "container_id"}

	// vms
	contVmsSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_containers_volume_mounts_%d", ts.UnixTimeStampNow()))
	contVmsSubCTE.TableName = "container_volume_mounts"
	contVmsSubCTE.Columns = []string{"volume_mount_id", "volume_mount_path", "volume_name"}
	contVmsRelationshipsSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_containers_volume_mounts_relationships_%d", ts.UnixTimeStampNow()))
	contVmsRelationshipsSubCTE.TableName = "containers_volume_mounts"
	contVmsRelationshipsSubCTE.Columns = []string{"chart_subcomponent_child_class_type_id", "container_id", "volume_mount_id"}

	// podSpec for containersMapByImageID
	podSpecSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_insert_spec_pod_template_containers_%d", ts.UnixTimeStampNow()))
	podSpecSubCTE.TableName = "chart_subcomponent_spec_pod_template_containers"
	podSpecSubCTE.Columns = []string{"chart_subcomponent_child_class_type_id", "container_id", "container_sort_order"}

	// probes
	contP := autogen_bases.ContainerProbes{}
	probesSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_container_probes_%d", ts.UnixTimeStampNow()))
	probesSubCTE.TableName = "container_probes"
	probesSubCTE.Columns = contP.GetTableColumns()
	conPRelationship := autogen_bases.ContainersProbes{}
	probesRelationshipsSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_containers_probes_relationship_%d", ts.UnixTimeStampNow()))
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
		contSubCTE.AddValues(cont.GetContainerID(), c.ContainerName, c.ContainerImageID, c.ContainerVersionTag, c.ContainerPlatformOs, c.ContainerRepository, c.ContainerImagePullPolicy, c.IsInitContainer)

		// pod spec to link container
		podSpecSubCTE.AddValues(p.GetPodSpecChildClassTypeID(), cont.GetContainerID(), sortOrderIndex)

		cmdArgsID := ts.UnixTimeStampNow()
		cmdArgsSubCTE.AddValues(cmdArgsID, cont.CmdArgs.CommandValues, cont.CmdArgs.ArgsValues)
		cmdArgsRelationshipsSubCTE.AddValues(cmdArgsID, cont.GetContainerID())

		computeResourcesID := ts.UnixTimeStampNow()
		computeResourcesSubCTE.AddValues(
			computeResourcesID,
			cont.ResourceRequest.ComputeResourcesCpuRequest,
			cont.ResourceRequest.ComputeResourcesCpuLimit,
			cont.ResourceRequest.ComputeResourcesRamRequest,
			cont.ResourceRequest.ComputeResourcesRamLimit,
			cont.ResourceRequest.ComputeResourcesEphemeralStorageRequest,
			cont.ResourceRequest.ComputeResourcesEphemeralStorageLimit,
		)
		computeResourcesRelationshipsSubCTE.AddValues(computeResourcesID, cont.GetContainerID())

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

			// cmdArgs
			cmdArgsSubCTE,
			cmdArgsRelationshipsSubCTE,

			//computeResources
			computeResourcesSubCTE,
			computeResourcesRelationshipsSubCTE,
		},
	}
	cteExpr.AppendSubCtes(templateMetadataCTE)
	cteExpr.AppendSubCtes(probeCTEs)
	cteExpr.AppendSubCtes(probeRelationshipCTEs)

	return cteExpr
}
