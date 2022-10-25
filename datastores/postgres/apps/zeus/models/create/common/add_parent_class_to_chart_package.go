package common

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func AddParentClassToChartPackage(c *create.Chart, parentClassID int) sql_query_templates.SubCTE {
	cpk := autogen_bases.ChartPackageComponents{
		ChartPackageID:                     c.GetChartPackageID(),
		ChartSubcomponentParentClassTypeID: parentClassID,
	}
	return createChartPackageComponentAddParentClassTypeSubCTE(&cpk)
}

func createChartPackageComponentAddParentClassTypeSubCTE(cpk *autogen_bases.ChartPackageComponents) sql_query_templates.SubCTE {
	var ts chronos.Chronos
	queryName := fmt.Sprintf("cte_%s_%d", cpk.GetTableName(), ts.UnixTimeStampNow())
	cpkAddParentClassTypeSubCTE := sql_query_templates.NewSubInsertCTE(queryName)
	cpkAddParentClassTypeSubCTE.TableName = cpk.GetTableName()
	cpkAddParentClassTypeSubCTE.Fields = cpk.GetTableColumns()
	cpkAddParentClassTypeSubCTE.Values = []apps.RowValues{cpk.GetRowValues(queryName)}
	return cpkAddParentClassTypeSubCTE
}
