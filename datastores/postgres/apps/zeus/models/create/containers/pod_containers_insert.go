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

	//
	contSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_insert_containers_%d", ts.UnixTimeStampNow()))
	contSubCTE.TableName = "containers"
	contSubCTE.Columns = []string{"container_id", "container_name", "container_image_id", "container_version_tag", "container_platform_os", "container_repository", "container_image_pull_policy", "is_init_container"}

	// ports
	portsSubCTE, portsRelationshipsSubCTE := CreatePortsCTEs()

	// env vars
	envVarsSubCTE, envVarsRelationshipsSubCTE := CreateEnvVarsCTEs()

	// cmd args
	cmdArgsSubCTE, cmdArgsRelationshipsSubCTE := CreateCmdArgsCTEs()

	// computeResources
	computeResourcesSubCTE, computeResourcesRelationshipsSubCTE := CreateComputeResourcesCTEs()

	// vms
	contVmsSubCTE, contVmsRelationshipsSubCTE := CreateVolumeMountsCTEs()

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

	// security context
	containerSecurityCtxSubCTE, containersSecurityCtxRelationSubCTE := CreateSecCtxCTEs()

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
		p.getContainerPortsValuesForInsert(cont, &portsSubCTE)
		p.getContainerPortsHeaderRelationshipValues(cont, &portsRelationshipsSubCTE)

		// env vars
		p.getInsertContainerEnvVarsValues(cont, &envVarsSubCTE)
		p.getContainerEnvVarRelationshipValues(cont, &envVarsRelationshipsSubCTE)

		// vms
		p.insertContainerVolumeMountsValues(cont, &contVmsSubCTE, &contVmsRelationshipsSubCTE)

		// security ctx
		p.insertContainerSecurityCtx(cont, &containerSecurityCtxSubCTE, &containersSecurityCtxRelationSubCTE)

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

			// security ctx
			containerSecurityCtxSubCTE,
			containersSecurityCtxRelationSubCTE,

			// cmdArgs
			cmdArgsSubCTE,
			cmdArgsRelationshipsSubCTE,

			//computeResources
			computeResourcesSubCTE,
			computeResourcesRelationshipsSubCTE,
		},
	}

	// insert new pod spec value share process value
	if p.Spec.ShareProcessNamespace != nil {
		shareNamespaceIDChildTypeID := ts.UnixTimeStampNow()
		podSpecShareNamespaceSubChildClassCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_podSpecShareProcessNamespaceCTE_%d", ts.UnixTimeStampNow()))
		podSpecShareNamespaceSubChildClassCTE.TableName = agCct.GetTableName()
		podSpecShareNamespaceSubChildClassCTE.Columns = []string{"chart_subcomponent_parent_class_type_id", "chart_subcomponent_child_class_type_id", "chart_subcomponent_child_class_type_name"}
		podSpecShareNamespaceSubChildClassCTE.AddValues(p.GetPodSpecParentClassTypeID(), shareNamespaceIDChildTypeID, "shareProcessNamespace")

		cv := autogen_bases.ChartSubcomponentsChildValues{}
		podSpecShareNamespaceSubChildClassValueCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_podSpecShareProcessNamespaceValueCTE_%d", ts.UnixTimeStampNow()))
		podSpecShareNamespaceSubChildClassValueCTE.TableName = cv.GetTableName()
		podSpecShareNamespaceSubChildClassValueCTE.Columns = []string{"chart_subcomponent_child_class_type_id", "chart_subcomponent_chart_package_template_injection", "chart_subcomponent_key_name", "chart_subcomponent_value"}
		podSpecShareNamespaceSubChildClassValueCTE.AddValues(shareNamespaceIDChildTypeID, p.Spec.ShareProcessNamespace.ChartSubcomponentChartPackageTemplateInjection, p.Spec.ShareProcessNamespace.ChartSubcomponentKeyName, p.Spec.ShareProcessNamespace.ChartSubcomponentValue)

		cteExpr.AppendSubCtes(sql_query_templates.SubCTEs{podSpecShareNamespaceSubChildClassCTE, podSpecShareNamespaceSubChildClassValueCTE})
	}

	cteExpr.AppendSubCtes(templateMetadataCTE)
	cteExpr.AppendSubCtes(probeCTEs)
	cteExpr.AppendSubCtes(probeRelationshipCTEs)

	return cteExpr
}
