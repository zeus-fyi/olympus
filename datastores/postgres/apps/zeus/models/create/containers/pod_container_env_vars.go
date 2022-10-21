package containers

import (
	"fmt"

	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
)

func (p *PodContainersGroup) insertContainerEnvVarsHeader() string {
	return "INSERT INTO container_environmental_vars(port_id, port_name, container_port, host_port) VALUES "
}

func (p *PodContainersGroup) getInsertContainerEnvVarsValues(parentExpression, containerImageID string) string {
	c, ok := p.Containers[containerImageID]
	if !ok {
		return ""
	}
	for _, ev := range c.Env {
		parentExpression += fmt.Sprintf("('%d', '%s', '%s')", ev.EnvID, ev.Name, ev.Value)
	}

	return parentExpression
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
