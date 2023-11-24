package jobs

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (j *Job) GetJobCTE(chart *charts.Chart) sql_query_templates.CTE {
	if chart != nil {
		j.SetChartPackageID(chart.GetChartPackageID())
	}
	var combinedSubCTEs sql_query_templates.SubCTEs
	chart.ChartComponentResourceID = JobChartComponentResourceID
	// metadata
	j.Metadata.Name.ChartSubcomponentValue = j.K8sJob.Name
	metaDataCtes := common.CreateParentMetadataSubCTEs(chart, j.Metadata)
	parentSpecCTE := common.CreateParentClassTypeSubCTE(chart, &j.Spec.ChartSubcomponentParentClassTypes)
	j.Spec.ChildClassSingleValue.SetParentClassTypeID(j.Spec.ParentClass.ChartSubcomponentParentClassTypeID)
	specCTE := common.CreateChildClassSingleValueSubCTEs(&j.Spec.ChildClassSingleValue)
	combinedSubCTEs = sql_query_templates.AppendSubCteSlices(metaDataCtes, parentSpecCTE, specCTE)
	cteExpr := sql_query_templates.CTE{
		Name:    "InsertJobCTEs",
		SubCTEs: combinedSubCTEs,
	}
	return cteExpr
}

func (j *Job) SetChartPackageID(id int) {
	j.Spec.ChartPackageID = id
	j.Metadata.ChartPackageID = id
}
