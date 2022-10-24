package common

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/common"
)

func InsertChildClassValues(parentExpression, parentClassTypeCteName string, cvals common.ChildClassMultiValue) string {
	cvTypeName := cvals.ChartSubcomponentChildClassTypes.ChartSubcomponentChildClassTypeName
	typeInsert, cvTypeNameCte := SetCvTypeInsert(cvTypeName, parentClassTypeCteName)

	valsToInsert := "VALUES "
	for i, val := range cvals.Values {
		valsToInsert += fmt.Sprintf("((SELECT chart_subcomponent_child_class_type_id FROM %s), '%s', '%s')", cvTypeNameCte, val.ChartSubcomponentKeyName, val.ChartSubcomponentValue)
		if i < len(cvals.Values)-1 {
			valsToInsert += ","
		}
	}
	cvChildCteName := fmt.Sprintf("cte_cvs_%s", cvTypeName)
	valueInsert := fmt.Sprintf(`
				%s AS (
					INSERT INTO chart_subcomponents_child_values(chart_subcomponent_child_class_type_id, chart_subcomponent_key_name, chart_subcomponent_value)
					%s
	),`, cvChildCteName, valsToInsert)
	returnExpression := fmt.Sprintf("%s %s %s", parentExpression, typeInsert, valueInsert)
	return returnExpression
}
