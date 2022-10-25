package containers

import (
	"encoding/json"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/containers"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (p *PodTemplateSpec) getInsertContainerEnvVarsValues(m map[string]containers.Container, imageID string, cteSubfield *sql_query_templates.SubCTE) {
	c, ok := m[imageID]
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

func (p *PodTemplateSpec) getContainerEnvVarRelationshipValues(m map[string]containers.Container, imageID string, cteSubfield *sql_query_templates.SubCTE) {
	c, ok := m[imageID]
	if !ok {
		return
	}

	podSpecChildClassTypeID := p.GetPodSpecChildClassTypeID()
	for _, ev := range c.GetEnvVars() {
		cteSubfield.AddValues(podSpecChildClassTypeID, c.GetContainerID(), ev.EnvID)
	}
	return
}
