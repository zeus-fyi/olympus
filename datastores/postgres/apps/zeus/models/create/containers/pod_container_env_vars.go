package containers

import (
	"encoding/json"
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/containers"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func CreateEnvVarsCTEs() (sql_query_templates.SubCTE, sql_query_templates.SubCTE) {
	// env vars
	ts := chronos.Chronos{}
	envVarsSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_container_environmental_vars_%d", ts.UnixTimeStampNow()))
	envVarsSubCTE.TableName = "container_environmental_vars"
	envVarsSubCTE.Columns = []string{"env_id", "name", "value"}
	envVarsRelationshipsSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_container_environmental_vars_relationships_%d", ts.UnixTimeStampNow()))
	envVarsRelationshipsSubCTE.TableName = "containers_environmental_vars"
	envVarsRelationshipsSubCTE.Columns = []string{"chart_subcomponent_child_class_type_id", "container_id", "env_id"}
	return envVarsSubCTE, envVarsRelationshipsSubCTE
}

func (p *PodTemplateSpec) getInsertContainerEnvVarsValues(c containers.Container, cteSubfield *sql_query_templates.SubCTE) {
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

func (p *PodTemplateSpec) getContainerEnvVarRelationshipValues(c containers.Container, cteSubfield *sql_query_templates.SubCTE) {

	podSpecChildClassTypeID := p.GetPodSpecChildClassTypeID()
	for _, ev := range c.GetEnvVars() {
		cteSubfield.AddValues(podSpecChildClassTypeID, c.GetContainerID(), ev.EnvID)
	}
	return
}
