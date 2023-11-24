package jobs

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (cj *CronJob) GetCronJobCTE(chart *charts.Chart) sql_query_templates.CTE {
	if chart != nil {
		cj.SetChartPackageID(chart.GetChartPackageID())
	}
	var combinedSubCTEs sql_query_templates.SubCTEs
	chart.ChartComponentResourceID = CronJobChartComponentResourceID
	// metadata
	cj.Metadata.Name.ChartSubcomponentValue = cj.K8sCronJob.Name
	metaDataCtes := common.CreateParentMetadataSubCTEs(chart, cj.Metadata)
	parentSpecCTE := common.CreateParentClassTypeSubCTE(chart, &cj.Spec.ChartSubcomponentParentClassTypes)
	cj.Spec.ChildClassSingleValue.SetParentClassTypeID(cj.Spec.ParentClass.ChartSubcomponentParentClassTypeID)
	specCTE := common.CreateChildClassSingleValueSubCTEs(&cj.Spec.ChildClassSingleValue)
	combinedSubCTEs = sql_query_templates.AppendSubCteSlices(metaDataCtes, parentSpecCTE, specCTE)
	cteExpr := sql_query_templates.CTE{
		Name:    "InsertCronJobCTEs",
		SubCTEs: combinedSubCTEs,
	}
	return cteExpr
}

func (cj *CronJob) SetChartPackageID(id int) {
	cj.Spec.ChartPackageID = id
	cj.Metadata.ChartPackageID = id
}
