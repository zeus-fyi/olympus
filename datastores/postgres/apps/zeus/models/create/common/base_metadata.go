package common

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

// CreateBaseMetadataSubCTEs returns parent cte, name, label, annotations sub ctes
func CreateBaseMetadataSubCTEs(metadata structs.Metadata) sql_query_templates.SubCTEs {
	var combinedSubCtes sql_query_templates.SubCTEs
	if metadata.HasName() {
		combinedSubCtes = sql_query_templates.AppendSubCteSlices(combinedSubCtes, CreateChildClassSingleValueSubCTEs(&metadata.Name))
	}
	if metadata.HasLabels() {
		combinedSubCtes = sql_query_templates.AppendSubCteSlices(combinedSubCtes, CreateChildClassMultiValueSubCTEs(&metadata.Labels))
	}
	if metadata.HasAnnotations() {
		combinedSubCtes = sql_query_templates.AppendSubCteSlices(combinedSubCtes, CreateChildClassMultiValueSubCTEs(&metadata.Annotations))
	}
	return combinedSubCtes
}
