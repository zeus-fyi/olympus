package containers

import (
	"encoding/json"
	"fmt"
)

func (p *PodContainersGroup) insertContainerEnvVarsHeader() string {
	return "INSERT INTO container_environmental_vars(env_id, name, value) VALUES "
}

func (p *PodContainersGroup) getInsertContainerEnvVarsValues(parentExpression, containerImageID string, isLastValuesGroup bool) string {
	c, ok := p.Containers[containerImageID]
	if !ok {
		return ""
	}
	for i, ev := range c.Env {

		jsonBvalue := "{}"
		if len(ev.Value) != 0 {
			bytes, _ := json.Marshal(ev.Value)
			jsonBvalue = string(bytes)
		}
		parentExpression += fmt.Sprintf("\n('%d', '%s', '%s'::jsonb)", ev.EnvID, ev.Name, jsonBvalue)
		if i < len(c.Env)-1 && !isLastValuesGroup {
			parentExpression += ","
		}
		i += 1
	}
	return parentExpression
}

func (p *PodContainersGroup) insertContainerEnvVarRelationshipHeader() string {
	return "INSERT INTO containers_environmental_vars(chart_subcomponent_child_class_type_id, container_id, env_id) VALUES "
}

func (p *PodContainersGroup) getContainerEnvVarRelationshipValues(parentExpression, containerImageID, classTypeID string, isLastValuesGroup bool) string {
	valsToInsert := ""
	c, ok := p.Containers[containerImageID]
	if !ok {
		return valsToInsert
	}
	for i, ev := range c.Env {
		parentExpression += fmt.Sprintf("\n('%s', (%s), '%d')", classTypeID, selectRelatedContainerIDFromImageID(containerImageID), ev.EnvID)
		if i < len(c.Env)-1 && !isLastValuesGroup {
			parentExpression += ","
		}
	}
	returnExpression := fmt.Sprintf("%s %s", parentExpression, valsToInsert)
	return returnExpression
}
