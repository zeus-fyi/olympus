package containers

import (
	"fmt"

	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/dev_hacks"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

// insertPodContainerGroupSQL, will use the next_id distributed ID generator and select the container id
// value for subsequent subcomponent relationships of its element, should greatly simplify the insert logic
func (p *PodContainersGroup) insertPodContainerGroupSQL(workloadChildGroupInfo autogen_bases.ChartSubcomponentChildClassTypes) string {
	i := len(p.Containers)

	// container
	contSubCTE := sql_query_templates.NewSubInsertCTE("cte_insert_containers")
	contSubCTE.TableName = "containers"
	contSubCTE.Fields = []string{"container_name", "container_image_id", "container_version_tag", "container_platform_os", "container_repository", "container_image_pull_policy"}

	// ports
	portsSubCTE := sql_query_templates.NewSubInsertCTE("cte_insert_container_ports")
	portsSubCTE.TableName = "container_ports"
	portsSubCTE.Fields = []string{"port_id", "port_name", "container_port", "host_port"}
	portsRelationshipsSubCTE := sql_query_templates.NewSubInsertCTE("cte_containers_ports_relationship")
	portsRelationshipsSubCTE.TableName = "containers_ports"
	portsSubCTE.Fields = []string{"chart_subcomponent_child_class_type_id", "container_id", "port_id"}

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
	contVmsSubCTE.Fields = []string{"chart_subcomponent_child_class_type_id", "container_id", "volume_mount_id"}

	cteExpr := sql_query_templates.CTE{
		Name: "insertPodContainerGroupSQL",
		SubCTEs: []sql_query_templates.SubCTE{
			contSubCTE,
			// ports
			portsSubCTE,
			portsRelationshipsSubCTE,
			// env vars
			envVarsSubCTE,
			envVarsRelationshipsSubCTE,
			// vms
			contVmsSubCTE,
			contVmsRelationshipsSubCTE,
		},
	}

	dev_hacks.Use(cteExpr)
	for _, _ = range p.Containers {

		// should continue appending values to header
		// container
		contSubCTE.AddValues()

		i += 1
	}

	fakeCte := " cte_term AS ( SELECT 1 ) SELECT true"

	return fakeCte
}

func (p *PodContainersGroup) generateHeaderIfNoneForCTE(parentExpression, header string) string {
	if len(parentExpression) <= 0 {
		parentExpression += header
	}
	return parentExpression
}

func (p *PodContainersGroup) getInsertContainerValues(c autogen_bases.Containers) string {
	processAndSetAmbiguousContainerFieldStatus(c)
	valsToInsert := fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', '%s')", c.ContainerName, c.ContainerImageID, c.ContainerVersionTag, c.ContainerPlatformOs, c.ContainerRepository, c.ContainerImagePullPolicy)
	return valsToInsert
}
