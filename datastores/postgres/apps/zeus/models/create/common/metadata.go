package common

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/common"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

// CreateMetadataSubCTEs returns name, label, annotations sub ctes
func CreateMetadataSubCTEs(metadata common.Metadata) (sql_query_templates.SubCTEs, sql_query_templates.SubCTEs, sql_query_templates.SubCTEs) {
	nameSubCtes := CreateChildClassSingleValueSubCTEs(metadata.Name)
	labelSubCtes := CreateChildClassMultiValueSubCTEs(metadata.Labels)
	annotationsSubCtes := CreateChildClassMultiValueSubCTEs(metadata.Annotations)
	return nameSubCtes, labelSubCtes, annotationsSubCtes
}
