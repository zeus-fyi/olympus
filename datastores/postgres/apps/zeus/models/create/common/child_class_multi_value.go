package common

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/common"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func CreateChildClassMultiValueSubCTEs(csv common.ChildClassMultiValue) sql_query_templates.SubCTEs {
	elementCount := len(csv.Values) + 1 // including child class type
	childClassMultiValueSubCTEs := make([]sql_query_templates.SubCTE, elementCount)
	for i, element := range csv.Values {
		if i == 0 {
			childClassTypeSubCTE := createChildClassSingleValueChildClassTypeSubCTE(csv.ChartSubcomponentChildClassTypes)
			childClassMultiValueSubCTEs[i] = childClassTypeSubCTE
		} else {
			childClassMultiValueSubCTEs[i] = createChildClassSingleValueSubCTE(element)
		}
	}
	return childClassMultiValueSubCTEs
}
