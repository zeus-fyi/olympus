package containers

import (
	"fmt"

	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
)

// insertPodContainerGroupSQL, will use the next_id distributed ID generator and select the container id
// value for subsequent subcomponent relationships of its element, should greatly simplify the insert logic
func (p *PodContainersGroup) insertPodContainerGroupSQL(workloadChildGroupInfo autogen_structs.ChartSubcomponentChildClassTypes) string {

	i := len(p.Containers)
	// header insert sql stmts
	// container
	containerInsertHeader := p.insertContainerHeader()

	// ports
	insertPortsParentExpressionHeader := p.insertContainerPortsHeader()
	insertContainerPortsHeaderRelationshipHeader := p.insertContainerPortsHeaderRelationshipHeader()
	// env vars
	insertContainerEnvVarsHeader := p.insertContainerEnvVarsHeader()
	insertContainerEnvVarRelationshipHeader := p.insertContainerEnvVarRelationshipHeader()

	// vms
	insertContainerVolumeMountsHeader := p.insertContainerVolumeMountsHeader()
	insertContainerVolumeMountRelationshipHeader := p.insertContainerVolumeMountRelationshipHeader()
	// should use imageID set when calling NewPodContainersGroupForDB
	for imageID, cont := range p.Containers {

		// should continue appending values to header
		// container
		containerInsertHeader = p.getInsertContainerValues(containerInsertHeader, cont.Metadata)

		// ports
		insertPortsParentExpressionHeader = p.getContainerPortsValuesForInsert(insertPortsParentExpressionHeader, imageID)
		classTypeID := fmt.Sprintf("%d", cont.ClassDefinition.ChartSubcomponentChildClassTypeID)
		insertContainerPortsHeaderRelationshipHeader = p.getContainerPortsHeaderRelationshipValues(insertContainerPortsHeaderRelationshipHeader, imageID, classTypeID)

		// env vars
		insertContainerEnvVarsHeader = p.getInsertContainerEnvVarsValues(insertContainerEnvVarsHeader, imageID)
		insertContainerEnvVarRelationshipHeader = p.getContainerEnvVarRelationshipValues(insertContainerEnvVarRelationshipHeader, imageID, classTypeID)

		// vol mounts
		insertContainerVolumeMountsHeader = p.getInsertContainerVolumeMountsValues(insertContainerVolumeMountsHeader, cont.VolumeMounts)
		insertContainerVolumeMountRelationshipHeader = p.getContainerVolumeMountRelationshipValues(insertContainerVolumeMountRelationshipHeader, imageID, classTypeID)
		if i < len(p.Containers)-1 {
			containerInsertHeader += ","
		}
		i += 1
	}

	q := fmt.Sprintf(
		`WITH cte_insert_containers AS (
					%s
				), cte_insert_container_ports AS (
				    %s
				), cte_containers_ports_relationship AS (
					%s
				), cte_container_environmental_vars AS (
					%s
				), cte_container_environmental_vars_relationships AS (
					%s
				), cte_containers_volume_mounts AS (
					%s
			    ) cte_containers_volume_mounts_relationships AS (
					%s 
				) 
		`,
		// containers
		containerInsertHeader,
		// ports
		insertPortsParentExpressionHeader, insertContainerPortsHeaderRelationshipHeader,
		// env vars
		insertContainerEnvVarsHeader, insertContainerEnvVarRelationshipHeader,
		// vm mounts
		insertContainerVolumeMountsHeader, insertContainerVolumeMountRelationshipHeader,
	)

	fakeCte := " cte_term AS ( SELECT 1 ) SELECT true"
	return q + fakeCte
}

func (p *PodContainersGroup) insertContainerHeader() string {
	containerInsert := "INSERT INTO containers(container_name, container_image_id, container_version_tag, container_platform_os, container_repository, container_image_pull_policy) VALUES \n"
	return containerInsert
}

func (p *PodContainersGroup) getInsertContainerValues(parentExpression string, c autogen_structs.Containers) string {
	processAndSetAmbiguousContainerFieldStatus(c)
	valsToInsert := fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', '%s')", c.ContainerName, c.ContainerImageID, c.ContainerVersionTag, c.ContainerPlatformOs, c.ContainerRepository, c.ContainerImagePullPolicy)
	returnExpression := fmt.Sprintf("%s \n %s", parentExpression, valsToInsert)
	return returnExpression
}
