package statefulsets

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (s *StatefulSet) GetStatefulSetCTE(chart *charts.Chart) sql_query_templates.CTE {
	var combinedSubCTEs sql_query_templates.SubCTEs
	chart.ChartComponentResourceID = StsChartComponentResourceID

	// metadata
	metaDataCtes := common.CreateParentMetadataSubCTEs(chart, s.Metadata)

	// sets spec as parent to sub elements
	s.SetSpecParentIDs()
	// spec
	specCtes := common.CreateSpecWorkloadTypeSubCTE(chart, s.Spec.SpecWorkload)

	// pod template spec
	podSpecTemplateCte := s.Spec.Template.InsertPodTemplateSpecContainersCTE(chart)
	podSpecTemplateSubCtes := podSpecTemplateCte.SubCTEs

	stsUpdateStrategyCTEs := common.CreateChildClassMultiValueSubCTEs(&s.Spec.StatefulSetUpdateStrategy)
	stsPodManagementPolicyCTEs := common.CreateChildClassSingleValueSubCTEs(&s.Spec.PodManagementPolicy)
	stsServiceNameCTEs := common.CreateChildClassSingleValueSubCTEs(&s.Spec.ServiceName)

	pvcCTEs := s.Spec.VolumeClaimTemplates.GetVCTemplateGroupSubCTEs(chart)

	combinedSubCTEs = sql_query_templates.AppendSubCteSlices(metaDataCtes, specCtes, podSpecTemplateSubCtes,
		stsUpdateStrategyCTEs, stsPodManagementPolicyCTEs, stsServiceNameCTEs, pvcCTEs)

	cteExpr := sql_query_templates.CTE{
		Name:    "InsertStatefulSetCTEs",
		SubCTEs: combinedSubCTEs,
	}
	return cteExpr
}
