package statefulset

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (s *StatefulSet) GetStatefulSetCTE(chart *charts.Chart) sql_query_templates.CTE {
	var combinedSubCTEs sql_query_templates.SubCTEs
	// metadata
	metaDataCtes := common.CreateParentMetadataSubCTEs(chart, s.Metadata)
	// spec
	specCtes := common.CreateSpecWorkloadTypeSubCTE(chart, s.Spec.SpecWorkload)

	// pod template spec
	podSpecTemplateCte := s.Spec.Template.InsertPodTemplateSpecContainersCTE(chart)
	podSpecTemplateSubCtes := podSpecTemplateCte.SubCTEs

	/*
		StatefulSetUpdateStrategy structs.ChildClassMultiValue
		PodManagementPolicy       structs.ChildClassSingleValue
		ServiceName               structs.ChildClassSingleValue
	*/
	stsUpdateStrategyCTEs := common.CreateChildClassMultiValueSubCTEs(&s.Spec.StatefulSetUpdateStrategy)
	stsPodManagementPolicyCTEs := common.CreateChildClassSingleValueSubCTEs(&s.Spec.PodManagementPolicy)
	stsServiceNameCTEs := common.CreateChildClassSingleValueSubCTEs(&s.Spec.ServiceName)

	combinedSubCTEs = sql_query_templates.AppendSubCteSlices(metaDataCtes, specCtes, podSpecTemplateSubCtes,
		stsUpdateStrategyCTEs, stsPodManagementPolicyCTEs, stsServiceNameCTEs)

	cteExpr := sql_query_templates.CTE{
		Name:    "InsertStatefulSetCTEs",
		SubCTEs: combinedSubCTEs,
	}
	return cteExpr
}
