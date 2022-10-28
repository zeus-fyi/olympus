package deployments

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (d *Deployment) GetDeploymentCTE(chart *charts.Chart) sql_query_templates.CTE {
	var combinedSubCTEs sql_query_templates.SubCTEs
	// metadata
	metaDataCtes := common.CreateParentMetadataSubCTEs(chart, d.Metadata)
	// spec
	specCtes := common.CreateSpecWorkloadTypeSubCTE(chart, d.Spec.SpecWorkload)

	// pod template spec
	podSpecTemplateCte := d.Spec.Template.InsertPodTemplateSpecContainersCTE(chart)
	podSpecTemplateSubCtes := podSpecTemplateCte.SubCTEs

	combinedSubCTEs = sql_query_templates.AppendSubCteSlices(metaDataCtes, specCtes, podSpecTemplateSubCtes)
	cteExpr := sql_query_templates.CTE{
		Name:    "InsertDeploymentCTEs",
		SubCTEs: combinedSubCTEs,
	}
	return cteExpr
}
