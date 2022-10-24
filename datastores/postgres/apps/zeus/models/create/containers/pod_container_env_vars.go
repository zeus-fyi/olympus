package containers

import (
	"encoding/json"

	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (p *PodContainersGroup) insertContainerEnvVarsHeader() string {
	return "INSERT INTO container_environmental_vars(env_id, name, value) VALUES "
}

func (p *PodContainersGroup) getInsertContainerEnvVarsValues(imageID string, cteSubfield *sql_query_templates.SubCTE) {
	c, ok := p.Containers[imageID]
	if !ok {
		return
	}
	for i, ev := range c.GetEnvVars() {
		jsonBvalue := "{}"
		if len(ev.Value) != 0 {
			bytes, _ := json.Marshal(ev.Value)
			jsonBvalue = string(bytes)
		}
		cteSubfield.AddValues(ev.EnvID, ev.Name, jsonBvalue)
		i += 1
	}
	return
}

func (p *PodContainersGroup) insertContainerEnvVarRelationshipHeader() string {
	return "INSERT INTO containers_environmental_vars(chart_subcomponent_child_class_type_id, container_id, env_id) VALUES "
}

func (p *PodContainersGroup) getContainerEnvVarRelationshipValues(podSpecChildClassTypeID int, imageID string, cteSubfield *sql_query_templates.SubCTE) {
	c, ok := p.Containers[imageID]
	if !ok {
		return
	}
	for _, ev := range c.GetEnvVars() {
		cteSubfield.AddValues(podSpecChildClassTypeID, selectRelatedContainerIDFromImageID(imageID), ev.EnvID)
	}
	return
}
