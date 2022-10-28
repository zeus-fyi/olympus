package ingresses

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (is *Spec) GetIngressSpecCTE(c *charts.Chart) sql_query_templates.SubCTEs {
	parentClassTypeSubCTE := common.CreateParentClassTypeSubCTE(c, &is.ChartSubcomponentParentClassTypes)
	pcID := is.ChartSubcomponentParentClassTypeID
	is.SetSpecParentIDs(pcID)
	chartComponentRelationshipCte := common.AddParentClassToChartPackage(c, pcID)
	// rules
	rulesCte := common.CreateSuperParentGroupClassTypeFromSlicesSubCTE(c, is.Rules.SuperParentClassGroup)
	// tls
	tlsCte := common.CreateSuperParentGroupClassTypeFromSlicesSubCTE(c, is.TLS.SuperParentClassGroup)
	combinedSubCtes := sql_query_templates.AppendSubCteSlices(parentClassTypeSubCTE, rulesCte, tlsCte, []sql_query_templates.SubCTE{chartComponentRelationshipCte})
	return combinedSubCtes
}

func (is *Spec) SetSpecParentIDs(id int) {
	is.ChartSubcomponentParentClassTypeID = id

	is.TLS.SetParentClassTypeIDs(id)
	is.Rules.SetParentClassTypeIDs(id)
}
