package common

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func CreateSuperParentGroupClassTypeChildrenFromSlicesSubCTE(spg structs.SuperParentClassGroup) sql_query_templates.SubCTEs {
	var childrenSubCTEs sql_query_templates.SubCTEs
	for _, pg := range spg.SuperParentClassSlice {
		if pg.ChildClassSingleValue != nil {
			csvSubCTEs := CreateChildClassSingleValueSubCTEs(pg.ChildClassSingleValue)
			childrenSubCTEs = sql_query_templates.AppendSubCteSlices(childrenSubCTEs, csvSubCTEs)
		}
		if pg.ChildClassMultiValue != nil {
			cmvSubCTEs := CreateChildClassMultiValueSubCTEs(pg.ChildClassMultiValue)
			childrenSubCTEs = sql_query_templates.AppendSubCteSlices(childrenSubCTEs, cmvSubCTEs)
		}
	}
	return sql_query_templates.AppendSubCteSlices(childrenSubCTEs)
}
