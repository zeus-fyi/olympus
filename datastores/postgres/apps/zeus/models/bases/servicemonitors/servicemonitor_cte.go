package servicemonitors

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (s *ServiceMonitor) GetServiceMonitorCTE(chart *charts.Chart) sql_query_templates.CTE {
	var combinedSubCTEs sql_query_templates.SubCTEs
	chart.ChartComponentResourceID = ServiceMonitorChartComponentResourceID
	// metadata
	metaDataCtes := common.CreateParentMetadataSubCTEs(chart, s.Metadata)
	s.SetSpecParentChildIDs()
	specCTE := common.CreateChildClassSingleValueSubCTEs(&s.Spec.ChildClassSingleValue)
	combinedSubCTEs = sql_query_templates.AppendSubCteSlices(metaDataCtes, specCTE)
	cteExpr := sql_query_templates.CTE{
		Name:    "InsertServiceMonitorCTEs",
		SubCTEs: combinedSubCTEs,
	}
	return cteExpr
}

func (s *ServiceMonitor) SetSpecParentChildIDs() {
	ts := chronos.Chronos{}
	parentID := ts.UnixTimeStampNow()
	childID := ts.UnixTimeStampNow()

	s.Spec.ChildClassSingleValue.SetParentClassTypeID(parentID)
	s.Spec.ChildClassSingleValue.SetChildClassTypeIDs(childID)
}
