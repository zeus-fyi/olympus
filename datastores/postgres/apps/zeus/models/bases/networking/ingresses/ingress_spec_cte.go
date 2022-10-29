package ingresses

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (is *Spec) GetIngressSpecCTE(c *charts.Chart) sql_query_templates.SubCTEs {
	parentClassTypeSubCTE := common.CreateParentClassTypeSubCTE(c, &is.ChartSubcomponentParentClassTypes)
	pcID := is.ChartSubcomponentParentClassTypeID
	is.SetSpecChartPackageResourceAndParentIDs(c.ChartPackageID, pcID)
	chartComponentRelationshipCte := common.AddParentClassToChartPackage(c, pcID)
	// rules
	rulesCte := common.CreateSuperParentGroupClassTypeChildrenFromSlicesSubCTE(is.Rules.SuperParentClassGroup)
	// tls
	tlsCte := common.CreateSuperParentGroupClassTypeChildrenFromSlicesSubCTE(is.TLS.SuperParentClassGroup)
	combinedSubCtes := sql_query_templates.AppendSubCteSlices(parentClassTypeSubCTE, rulesCte, tlsCte)
	if is.IngressClassName != nil {
		combinedSubCtes = sql_query_templates.AppendSubCteSlices(combinedSubCtes, common.CreateChildClassSingleValueSubCTEs(&is.IngressClassName.ChildClassSingleValue))
	}
	return sql_query_templates.AppendSubCteSlices(combinedSubCtes, []sql_query_templates.SubCTE{chartComponentRelationshipCte})
}
