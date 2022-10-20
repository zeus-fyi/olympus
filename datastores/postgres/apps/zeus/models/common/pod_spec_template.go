package common

import (
	"fmt"

	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
)

func SetPodSpecTemplateChildTypeInsert(parentExpression, parentClassTypeCteName string, psChildClassType autogen_structs.ChartSubcomponentChildClassTypes) string {
	cvTypeName := psChildClassType.ChartSubcomponentChildClassTypeName
	typeInsert, cvTypeNameCte := SetCvTypeInsert(cvTypeName, parentClassTypeCteName)

	// TODO needs to link containers into here, so probably even create the containers completely independently
	cvChildCteName := fmt.Sprintf("cte_ps_cvs_%s", cvTypeName)
	valueInsert := fmt.Sprintf(`
				%s AS (
					INSERT INTO chart_subcomponent_spec_pod_template_containers(chart_subcomponent_child_class_type_id, container_id, is_init_container)
					VALUES ((SELECT chart_subcomponent_child_class_type_id FROM %s), '%s', '%s')
	),`, cvChildCteName, cvTypeNameCte, "", "")

	returnExpression := fmt.Sprintf("%s %s %s", parentExpression, typeInsert, valueInsert)
	return returnExpression
}
