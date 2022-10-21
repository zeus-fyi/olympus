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
		cont.ProcessAndSetAmbiguousContainerFieldStatusAndSubfieldIds()
		containerInsertHeader = p.getInsertContainerValues(containerInsertHeader, cont.Metadata)

		shouldAppendCommaBool := i < len(p.Containers)-1
		if shouldAppendCommaBool {
			containerInsertHeader += ","
		}
		// ports
		insertPortsParentExpressionHeader = p.getContainerPortsValuesForInsert(insertPortsParentExpressionHeader, imageID, shouldAppendCommaBool)
		classTypeID := fmt.Sprintf("%d", cont.ClassDefinition.ChartSubcomponentChildClassTypeID)
		insertContainerPortsHeaderRelationshipHeader = p.getContainerPortsHeaderRelationshipValues(insertContainerPortsHeaderRelationshipHeader, imageID, classTypeID)

		// env vars
		insertContainerEnvVarsHeader = p.getInsertContainerEnvVarsValues(insertContainerEnvVarsHeader, imageID, shouldAppendCommaBool)
		insertContainerEnvVarRelationshipHeader = p.getContainerEnvVarRelationshipValues(insertContainerEnvVarRelationshipHeader, imageID, classTypeID, shouldAppendCommaBool)

		// vol mounts
		insertContainerVolumeMountsHeader = p.getInsertContainerVolumeMountsValues(insertContainerVolumeMountsHeader, cont.VolumeMounts, shouldAppendCommaBool)
		insertContainerVolumeMountRelationshipHeader = p.getContainerVolumeMountRelationshipValues(insertContainerVolumeMountRelationshipHeader, imageID, classTypeID)

		i += 1
	}

	// env vars
	cteEnvVars := AddCTEIfValuesExist(insertContainerEnvVarsHeader, p.insertContainerEnvVarsHeader(), "cte_container_environmental_vars")
	cteEnvVarsRelationships := AddCTEIfValuesExist(insertContainerEnvVarRelationshipHeader, p.insertContainerEnvVarRelationshipHeader(), "cte_container_environmental_vars_relationships")

	// vm mounts
	cteContainerVms := AddCTEIfValuesExist(insertContainerVolumeMountsHeader, p.insertContainerVolumeMountsHeader(), "cte_containers_volume_mounts")
	cteContainerVmsRelationships := AddCTEIfValuesExist(insertContainerVolumeMountsHeader, p.insertContainerVolumeMountRelationshipHeader(), "cte_containers_volume_mounts_relationships")

	// ports
	ctePorts := AddCTEIfValuesExist(insertPortsParentExpressionHeader, p.insertContainerPortsHeader(), "cte_insert_container_ports")
	ctePortsRelationships := AddCTEIfValuesExist(insertContainerPortsHeaderRelationshipHeader, p.insertContainerPortsHeaderRelationshipHeader(), "cte_containers_ports_relationship")

	q := fmt.Sprintf(
		`WITH cte_insert_containers AS (
					%s
				),  %s %s %s %s %s %s `,
		// containers
		containerInsertHeader,
		// ports
		ctePorts, ctePortsRelationships,
		// env vars
		cteEnvVars, cteEnvVarsRelationships,
		// vm mounts
		cteContainerVms, cteContainerVmsRelationships,
	)

	fakeCte := " cte_term AS ( SELECT 1 ) SELECT true"
	q += fakeCte
	return q
}

func AddCTEIfValuesExist(parentExpression, startingHeader, cteHeader string) string {
	if len(parentExpression) > len(startingHeader) {
		return fmt.Sprintf("%s AS (\n\t%s\n), ", cteHeader, parentExpression)
	}
	return ""
}

func (p *PodContainersGroup) generateHeaderIfNoneForCTE(parentExpression, header string) string {
	if len(parentExpression) <= 0 {
		parentExpression += header
	}
	return parentExpression
}

func (p *PodContainersGroup) insertContainerHeader() string {
	containerInsert := "INSERT INTO containers(container_name, container_image_id, container_version_tag, container_platform_os, container_repository, container_image_pull_policy) VALUES "
	return containerInsert
}

func (p *PodContainersGroup) getInsertContainerValues(parentExpression string, c autogen_structs.Containers) string {
	processAndSetAmbiguousContainerFieldStatus(c)
	valsToInsert := fmt.Sprintf("\n ('%s', '%s', '%s', '%s', '%s', '%s')", c.ContainerName, c.ContainerImageID, c.ContainerVersionTag, c.ContainerPlatformOs, c.ContainerRepository, c.ContainerImagePullPolicy)
	returnExpression := fmt.Sprintf("%s %s", parentExpression, valsToInsert)
	return returnExpression
}
