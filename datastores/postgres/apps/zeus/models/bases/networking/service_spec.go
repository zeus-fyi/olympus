package networking

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

// ServiceSpec has these type options: ClusterIP, NodePort, LoadBalancer, ExternalName
type ServiceSpec struct {
	autogen_bases.ChartSubcomponentParentClassTypes
	Type     structs.ChildClassSingleValue
	Selector structs.Selector
	ServicePorts
}

func NewServiceSpec() ServiceSpec {
	s := ServiceSpec{}
	s.ChartSubcomponentParentClassTypeName = "ServiceSpec"
	s.ChartComponentResourceID = SvcChartComponentResourceID
	s.Type = structs.NewChildClassSingleValue("type")
	s.ServicePorts = NewServicePorts()
	s.Selector = structs.NewSelector()
	return s
}

func (ss *ServiceSpec) CreateServiceSpecSubCTE(c *charts.Chart) sql_query_templates.SubCTEs {
	parentClassTypeSubCTE := common.CreateParentClassTypeSubCTE(c, &ss.ChartSubcomponentParentClassTypes)
	pcID := ss.ChartSubcomponentParentClassTypeID
	ss.SetParentIDs(pcID)
	chartComponentRelationshipCte := common.AddParentClassToChartPackage(c, pcID)
	matchLabelsCtes := common.CreateChildClassMultiValueSubCTEs(&ss.Selector.MatchLabels)
	portsCte := common.CreateFromSliceChildClassMultiValueSubCTEs(ss.Ports)
	combinedSubCtes := sql_query_templates.AppendSubCteSlices(parentClassTypeSubCTE, matchLabelsCtes, portsCte, []sql_query_templates.SubCTE{chartComponentRelationshipCte})
	return combinedSubCtes
}

func (ss *ServiceSpec) SetParentIDs(id int) {
	ss.ChartSubcomponentParentClassTypeID = id
	ss.Type.ChartSubcomponentParentClassTypeID = id
	ss.Selector.MatchLabels.ChartSubcomponentParentClassTypeID = id

	for i, _ := range ss.Ports {
		ss.Ports[i].ChartSubcomponentParentClassTypeID = id
	}

}
