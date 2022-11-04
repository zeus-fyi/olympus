package structs

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (cm *ChildMetadata) GetSubCTEs() sql_query_templates.SubCTEs {

	nameCTEs := common.CreateChildClassSingleValueSubCTEs(&cm.Name)
	annotationCTEs := common.CreateChildClassMultiValueSubCTEs(&cm.Annotations)
	labelsCTEs := common.CreateChildClassMultiValueSubCTEs(&cm.Labels)
	return sql_query_templates.AppendSubCteSlices(nameCTEs, annotationCTEs, labelsCTEs)
}
