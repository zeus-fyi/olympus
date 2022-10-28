package ingresses

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (i *Ingress) GetIngressCTE(chart *charts.Chart) sql_query_templates.CTE {
	if chart != nil {
		i.SetChartPackageID(chart.GetChartPackageID())
	}
	var combinedSubCTEs sql_query_templates.SubCTEs
	// metadata
	metaDataCtes := common.CreateParentMetadataSubCTEs(chart, i.Metadata)
	// spec
	specCtes := i.GetIngressSpecCTE(chart)
	combinedSubCTEs = sql_query_templates.AppendSubCteSlices(metaDataCtes, specCtes)
	cteExpr := sql_query_templates.CTE{
		Name:    "InsertIngressCTEs",
		SubCTEs: combinedSubCTEs,
	}
	return cteExpr
}
