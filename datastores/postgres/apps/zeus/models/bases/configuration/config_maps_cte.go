package configuration

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (cm *ConfigMap) GetConfigMapCTE(chart *charts.Chart) sql_query_templates.CTE {
	if chart != nil {
		cm.SetChartPackageID(chart.GetChartPackageID())
	}
	var combinedSubCTEs sql_query_templates.SubCTEs
	// metadata
	metaDataCtes := common.CreateParentMetadataSubCTEs(chart, cm.Metadata)
	// spec
	dataCtes := cm.GetConfigMapDataCTE(chart)
	combinedSubCTEs = sql_query_templates.AppendSubCteSlices(metaDataCtes, dataCtes)
	cteExpr := sql_query_templates.CTE{
		Name:    "InsertConfigMapCTEs",
		SubCTEs: combinedSubCTEs,
	}
	return cteExpr
}

func (cm *ConfigMap) GetConfigMapDataCTE(chart *charts.Chart) sql_query_templates.SubCTEs {
	if chart != nil {
		cm.SetChartPackageID(chart.GetChartPackageID())
		cm.Data.ChartPackageID = chart.GetChartPackageID()
		cm.Data.ChartComponentResourceID = ConfigMapChartComponentResourceID
	}
	spgSubCTEs := common.CreateSuperParentClassTypeParentChildrenAndChartPackageRelationSubCTE(chart, cm.Data.SuperParentClass)
	return spgSubCTEs

}
