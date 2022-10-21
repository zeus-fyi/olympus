package containers

import (
	"fmt"

	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
)

func selectRelatedContainerIDFromImageID(imageID string) string {
	selectRelatedContainerID := fmt.Sprintf("SELECT container_id FROM containers WHERE container_image_id = %s", imageID)
	return selectRelatedContainerID
}

// insertPodContainerGroupSQL, will use the next_id distributed ID generator and select the container id
// value for subsequent subcomponent relationships of its element, should greatly simplify the insert logic
func (p *PodContainersGroup) insertPodContainerGroupSQL(workloadChildGroupInfo autogen_structs.ChartSubcomponentChildClassTypes) string {
	valsToInsert := "VALUES "

	i := len(p.Containers)
	insertPortsParentExpression := p.insertContainerPortsHeader()
	// should use imageID set when calling NewPodContainersGroupForDB
	for imageID, cont := range p.Containers {
		c := cont.Metadata
		processAndSetAmbiguousContainerFieldStatus(c)
		valsToInsert += fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', '%s')", c.ContainerName, c.ContainerImageID, c.ContainerVersionTag, c.ContainerPlatformOs, c.ContainerRepository, c.ContainerImagePullPolicy)
		if i < len(p.Containers)-1 {
			valsToInsert += ","
		}

		// should continue appending values to header
		insertPortsParentExpression = p.getContainerPortsValuesForInsert(insertPortsParentExpression, imageID)

		i += 1
	}

	containerInsert := fmt.Sprintf(`
				%s AS (
					INSERT INTO containers(container_name, container_image_id, container_version_tag, container_platform_os, container_repository, container_image_pull_policy)
					%s
	),`, "cte_containers", valsToInsert)

	q := fmt.Sprintf(
		`WITH cte_insert_containers AS (
					%s
				), cte_insert_container_ports AS (
				    %s
				), cte_container_environmental_vars AS (
				), cte_compute_resources_key_values_jsonb AS (
			    ), 
		`, containerInsert, insertPortsParentExpression,
	)
	return q
}
