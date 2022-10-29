package common

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func CreateChildClassSingleValueSubCTEs(csv *structs.ChildClassSingleValue) sql_query_templates.SubCTEs {
	if csv.GetChildClassTypeID() == 0 || csv.GetChildValueTypeID() == 0 {
		var ts chronos.Chronos
		classTypeID := ts.UnixTimeStampNow()
		csv.SetChildClassTypeIDs(classTypeID)
	}
	childClassTypeSubCTE := createChildClassSingleValueChildClassTypeSubCTE(&csv.ChartSubcomponentChildClassTypes)
	childClassTypeName := csv.ChartSubcomponentChildClassTypeName
	childClassValueTypeSubCTE := createChildClassSingleValueSubCTE(childClassTypeName, &csv.ChartSubcomponentsChildValues)
	return []sql_query_templates.SubCTE{childClassTypeSubCTE, childClassValueTypeSubCTE}
}

func createChildClassSingleValueSubCTE(childClassTypeName string, csv *autogen_bases.ChartSubcomponentsChildValues) sql_query_templates.SubCTE {
	var ts chronos.Chronos
	queryName := fmt.Sprintf("cte_%s_value_%d", childClassTypeName, ts.UnixTimeStampNow())
	subCTE := sql_query_templates.NewSubInsertCTE(queryName)
	subCTE.TableName = csv.GetTableName()
	subCTE.Columns = csv.GetTableColumns()
	subCTE.Values = []apps.RowValues{csv.GetRowValues(queryName)}
	return subCTE
}

func createChildClassSingleValueChildClassTypeSubCTE(csvType *autogen_bases.ChartSubcomponentChildClassTypes) sql_query_templates.SubCTE {
	childClassTypeName := csvType.ChartSubcomponentChildClassTypeName
	var ts chronos.Chronos
	queryName := fmt.Sprintf("cte_%s_%d", childClassTypeName, ts.UnixTimeStampNow())
	childClassTypeSubCTE := sql_query_templates.NewSubInsertCTE(queryName)
	childClassTypeSubCTE.TableName = csvType.GetTableName()
	childClassTypeSubCTE.Columns = csvType.GetTableColumns()
	childClassTypeSubCTE.Values = []apps.RowValues{csvType.GetRowValues(queryName)}
	return childClassTypeSubCTE
}
