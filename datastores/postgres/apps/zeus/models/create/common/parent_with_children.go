package common

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
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

func CreateSuperParentClassTypeParentChildrenAndChartPackageRelationSubCTE(c *charts.Chart, pg structs.SuperParentClass) sql_query_templates.SubCTEs {
	pcID := pg.GetSuperParentClassTypeID()
	pcSubCte := AddParentClassToChartPackage(c, pcID)
	parentDataCTE := CreateParentClassTypeSubCTE(c, &pg.ChartSubcomponentParentClassTypes)

	var childrenSubCTEs sql_query_templates.SubCTEs
	if pg.ChildClassSingleValue != nil {
		csvSubCTEs := CreateChildClassSingleValueSubCTEs(pg.ChildClassSingleValue)
		childrenSubCTEs = sql_query_templates.AppendSubCteSlices(childrenSubCTEs, csvSubCTEs)
	}
	if pg.ChildClassMultiValue != nil {
		cmvSubCTEs := CreateChildClassMultiValueSubCTEs(pg.ChildClassMultiValue)
		childrenSubCTEs = sql_query_templates.AppendSubCteSlices(childrenSubCTEs, cmvSubCTEs)
	}
	return sql_query_templates.AppendSubCteSlices(childrenSubCTEs, parentDataCTE, []sql_query_templates.SubCTE{pcSubCte})
}
