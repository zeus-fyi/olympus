package common

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

// CreateMetadataSubCTEs returns child type cte, name, label, annotations sub ctes
func CreateMetadataSubCTEs(metadata structs.Metadata) sql_query_templates.SubCTEs {
	var combinedCTEs sql_query_templates.SubCTEs
	if metadata.HasName() {
		combinedCTEs = sql_query_templates.AppendSubCteSlices(combinedCTEs, CreateChildClassSingleValueSubCTEs(&metadata.Name))
	}
	if metadata.HasLabels() {
		combinedCTEs = sql_query_templates.AppendSubCteSlices(combinedCTEs, CreateChildClassMultiValueSubCTEs(&metadata.Labels))
	}
	if metadata.HasAnnotations() {
		combinedCTEs = sql_query_templates.AppendSubCteSlices(combinedCTEs, CreateChildClassMultiValueSubCTEs(&metadata.Annotations))
	}
	return combinedCTEs
}
