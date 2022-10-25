package common

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

// CreateParentMetadataSubCTEs returns parent cte, name, label, annotations sub ctes
func CreateParentMetadataSubCTEs(c *create.Chart, metadata structs.ParentMetaData) sql_query_templates.SubCTEs {
	if metadata.ChartSubcomponentParentClassTypeID == 0 {
		ts := chronos.Chronos{}
		metadata.SetParentClassTypeIDs(ts.UnixTimeStampNow())
	}

	parentSubCte := CreateParentClassTypeSubCTE(&metadata.ChartSubcomponentParentClassTypes)
	nameSubCtes := CreateChildClassSingleValueSubCTEs(&metadata.Name)
	labelSubCtes := CreateChildClassMultiValueSubCTEs(&metadata.Labels)
	annotationsSubCtes := CreateChildClassMultiValueSubCTEs(&metadata.Annotations)
	chartComponentRelationship := AddParentClassToChartPackage(c, metadata.ChartSubcomponentParentClassTypeID)
	combinedSubCtes := sql_query_templates.AppendSubCteSlices(parentSubCte, nameSubCtes, labelSubCtes, annotationsSubCtes, []sql_query_templates.SubCTE{chartComponentRelationship})
	return combinedSubCtes
}
