package common

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/common"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func CreateChildClassSingleValueSubCTEs(csv common.ChildClassSingleValue) sql_query_templates.SubCTEs {
	childClassTypeSubCTE := createChildClassSingleValueChildClassTypeSubCTE(csv.ChartSubcomponentChildClassTypes)
	childClassValueTypeSubCTE := createChildClassSingleValueSubCTE(csv.ChartSubcomponentsChildValues)
	return []sql_query_templates.SubCTE{childClassTypeSubCTE, childClassValueTypeSubCTE}
}

func createChildClassSingleValueSubCTE(csv autogen_bases.ChartSubcomponentsChildValues) sql_query_templates.SubCTE {
	keyName := csv.ChartSubcomponentKeyName
	subCTE := sql_query_templates.NewSubInsertCTE("cte_" + keyName)
	subCTE.TableName = "chart_subcomponents_child_values"
	subCTE.Fields = []string{"chart_subcomponent_child_class_type_id", "chart_subcomponent_key_name", "chart_subcomponent_value"}
	return subCTE
}

func createChildClassSingleValueChildClassTypeSubCTE(csv autogen_bases.ChartSubcomponentChildClassTypes) sql_query_templates.SubCTE {
	childClassType := csv.ChartSubcomponentChildClassTypeName
	childClassTypeSubCTE := sql_query_templates.NewSubInsertCTE("cte_" + childClassType)
	childClassTypeSubCTE.TableName = "chart_subcomponent_child_class_types"
	childClassTypeSubCTE.Fields = []string{"chart_subcomponent_parent_class_type_id", "chart_subcomponent_child_class_type_name"}
	return childClassTypeSubCTE
}
