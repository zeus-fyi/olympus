package containers

import (
	"fmt"

	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/containers"
)

func (p *PodContainersGroup) insertContainerEnvVars(parentExpression, containerImageID string, envVars containers.ContainerEnvVars, workloadChildGroupInfo autogen_structs.ChartSubcomponentChildClassTypes) string {
	valsToInsert := "VALUES "

	for i, ev := range envVars {
		valsToInsert += fmt.Sprintf("('%d', '%s', '%s')", ev.EnvID, ev.Name, ev.Value)
		if i < len(envVars)-1 {
			valsToInsert += ","
		}
	}

	containerInsert := fmt.Sprintf(`
				%s AS (
					INSERT INTO container_environmental_vars(env_id, name, value)
					%s
	),`, "cte_container_environmental_vars", valsToInsert)

	returnExpression := fmt.Sprintf("%s %s", parentExpression, containerInsert)
	return returnExpression
}

func (p *PodContainersGroup) insertContainerEnvVarRelationship(parentExpression, containerImageID string, envVar autogen_structs.ContainerEnvironmentalVars, cct autogen_structs.ChartSubcomponentChildClassTypes) string {
	valsToInsert := "VALUES "
	valsToInsert += fmt.Sprintf("('%d', (%s), '%d')", cct.ChartSubcomponentChildClassTypeID, selectRelatedContainerIDFromImageID(containerImageID), envVar.EnvID)
	containerInsert := fmt.Sprintf(`
				%s AS (
					INSERT INTO containers_environmental_vars(chart_subcomponent_child_class_type_id, container_id, env_id)
					%s
	),`, "cte_containers_environmental_vars", valsToInsert)

	returnExpression := fmt.Sprintf("%s %s", parentExpression, containerInsert)
	return returnExpression
}
