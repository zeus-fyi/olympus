package common

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/common"
)

func InsertChildClassSingleValueType(parentExpression, parentClassTypeCteName string, sv common.ChildClassSingleValue) string {
	cvTypeName := sv.ChartSubcomponentChildClassTypes.ChartSubcomponentChildClassTypeName
	typeInsert, cvTypeNameCte := SetCvTypeInsert(cvTypeName, parentClassTypeCteName)
	cvChildCteName := fmt.Sprintf("cte_cv_%s", cvTypeName)
	valueInsert := fmt.Sprintf(`
				%s AS (
					INSERT INTO chart_subcomponents_child_values(chart_subcomponent_child_class_type_id, chart_subcomponent_key_name, chart_subcomponent_value)
					VALUES ((SELECT chart_subcomponent_child_class_type_id FROM %s), '%s', '%s')
	),`, cvChildCteName, cvTypeNameCte, sv.ChartSubcomponentKeyName, sv.ChartSubcomponentValue)

	returnExpression := fmt.Sprintf("%s %s %s", parentExpression, typeInsert, valueInsert)
	return returnExpression
}

func SetCvTypeInsert(cvTypeName, parentClassTypeCteName string) (string, string) {
	cvTypeNameCte := fmt.Sprintf("cte_%s", cvTypeName)
	typeInsert := fmt.Sprintf(`
				%s AS (
					INSERT INTO chart_subcomponent_child_class_types(chart_subcomponent_parent_class_type_id, chart_subcomponent_child_class_type_name)
					VALUES ((SELECT chart_subcomponent_parent_class_type_id FROM %s), '%s')
					RETURNING chart_subcomponent_child_class_type_id
	),`, cvTypeNameCte, parentClassTypeCteName, cvTypeName)
	return typeInsert, cvTypeNameCte
}
