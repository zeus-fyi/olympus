package common

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func CreateSuperParentGroupClassTypeFromSlicesSubCTE(c *charts.Chart, spg structs.SuperParentClassGroup) sql_query_templates.SubCTEs {
	pcTypeSubCTEs := CreateParentClassTypeSubCTE(c, &spg.ChartSubcomponentParentClassTypes)
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
	return sql_query_templates.AppendSubCteSlices(pcTypeSubCTEs, childrenSubCTEs)
}
