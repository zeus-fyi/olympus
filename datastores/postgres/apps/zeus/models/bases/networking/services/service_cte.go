package services

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (s *Service) GetServiceCTE(chart *charts.Chart) sql_query_templates.CTE {
	if chart != nil {
		s.SetChartPackageID(chart.GetChartPackageID())
	}
	var combinedSubCTEs sql_query_templates.SubCTEs
	// metadata
	metaDataCtes := common.CreateParentMetadataSubCTEs(chart, s.Metadata)
	// spec
	specCtes := s.CreateServiceSpecSubCTE(chart)
	combinedSubCTEs = sql_query_templates.AppendSubCteSlices(metaDataCtes, specCtes)
	cteExpr := sql_query_templates.CTE{
		Name:    "InsertServiceCTEs",
		SubCTEs: combinedSubCTEs,
	}
	return cteExpr
}
