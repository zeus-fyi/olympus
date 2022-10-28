package common

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

// CreateParentMetadataSubCTEs returns parent cte, name, label, annotations sub ctes
func CreateParentMetadataSubCTEs(c *charts.Chart, metadata structs.ParentMetaData) sql_query_templates.SubCTEs {
	if metadata.ChartSubcomponentParentClassTypeID == 0 {
		ts := chronos.Chronos{}
		metadata.SetParentClassTypeIDs(ts.UnixTimeStampNow())
	}
	parentSubCte := CreateParentClassTypeSubCTE(c, &metadata.ChartSubcomponentParentClassTypes)
	if metadata.HasName() {
		parentSubCte = sql_query_templates.AppendSubCteSlices(parentSubCte, CreateChildClassSingleValueSubCTEs(&metadata.Name))
	}
	if metadata.HasLabels() {
		parentSubCte = sql_query_templates.AppendSubCteSlices(parentSubCte, CreateChildClassMultiValueSubCTEs(&metadata.Labels))
	}
	if metadata.HasAnnotations() {
		parentSubCte = sql_query_templates.AppendSubCteSlices(parentSubCte, CreateChildClassMultiValueSubCTEs(&metadata.Annotations))
	}
	chartComponentRelationship := AddParentClassToChartPackage(c, metadata.ChartSubcomponentParentClassTypeID)
	combinedSubCtes := sql_query_templates.AppendSubCteSlices(parentSubCte, []sql_query_templates.SubCTE{chartComponentRelationship})
	return combinedSubCtes
}
