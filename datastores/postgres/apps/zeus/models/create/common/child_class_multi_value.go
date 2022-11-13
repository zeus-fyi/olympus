package common

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func CreateChildClassMultiValueSubCTEs(csv *structs.ChildClassMultiValue) sql_query_templates.SubCTEs {
	if csv.GetMultiValueChildClassTypeID() == 0 {
		var ts chronos.Chronos
		classTypeID := ts.UnixTimeStampNow()
		csv.SetChildClassTypeIDs(classTypeID)
	}
	childClassTypeSubCTE := createChildClassSingleValueChildClassTypeSubCTE(&csv.ChartSubcomponentChildClassTypes)
	childClassTypeSubCTESlice := []sql_query_templates.SubCTE{childClassTypeSubCTE}
	childClassTypeName := csv.ChartSubcomponentChildClassTypeName
	childClassValuesSubCTEs := make([]sql_query_templates.SubCTE, len(csv.Values))
	for i, element := range csv.Values {
		cteName := childClassTypeName + fmt.Sprintf("_%d", i)
		childClassValuesSubCTEs[i] = createChildClassSingleValueSubCTE(cteName, &element)
	}
	return sql_query_templates.AppendSubCteSlices(childClassTypeSubCTESlice, childClassValuesSubCTEs)
}

func CreateFromSliceChildClassMultiValueSubCTEs(csvSlice []structs.ChildClassMultiValue) sql_query_templates.SubCTEs {
	var combinedCtes sql_query_templates.SubCTEs
	for _, csv := range csvSlice {
		newChildValueCtes := CreateChildClassMultiValueSubCTEs(&csv)
		combinedCtes = sql_query_templates.AppendSubCteSlices(combinedCtes, newChildValueCtes)
	}

	return combinedCtes
}
