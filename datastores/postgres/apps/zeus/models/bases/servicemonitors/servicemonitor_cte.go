package servicemonitors

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (s *ServiceMonitor) GetServiceMonitorCTE(chart *charts.Chart) sql_query_templates.CTE {
	if chart != nil {
		s.SetChartPackageID(chart.GetChartPackageID())
	}
	var combinedSubCTEs sql_query_templates.SubCTEs
	chart.ChartComponentResourceID = ServiceMonitorChartComponentResourceID
	// metadata
	metaDataCtes := common.CreateParentMetadataSubCTEs(chart, s.Metadata)
	parentSpecCTE := common.CreateParentClassTypeSubCTE(chart, &s.Spec.ChartSubcomponentParentClassTypes)
	s.Spec.ChildClassSingleValue.SetParentClassTypeID(s.Spec.ParentClass.ChartSubcomponentParentClassTypeID)
	specCTE := common.CreateChildClassSingleValueSubCTEs(&s.Spec.ChildClassSingleValue)
	combinedSubCTEs = sql_query_templates.AppendSubCteSlices(metaDataCtes, parentSpecCTE, specCTE)
	cteExpr := sql_query_templates.CTE{
		Name:    "InsertServiceMonitorCTEs",
		SubCTEs: combinedSubCTEs,
	}
	return cteExpr
}

func (s *ServiceMonitor) SetChartPackageID(id int) {
	s.Spec.ChartPackageID = id
	s.Metadata.ChartPackageID = id
}
