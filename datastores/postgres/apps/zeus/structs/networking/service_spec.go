package networking

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/common"
)

// ServiceSpec has these type options: ClusterIP, NodePort, LoadBalancer, ExternalName
type ServiceSpec struct {
	Type     autogen_structs.ChartSubcomponentsChildValues
	Selector common.Selector
	Ports    ServicePorts
}

func NewServiceSpec() ServiceSpec {
	s := ServiceSpec{}
	s.Type = autogen_structs.ChartSubcomponentsChildValues{
		ChartSubcomponentChildClassTypeID:              0,
		ChartSubcomponentChartPackageTemplateInjection: true,
		ChartSubcomponentKeyName:                       "type",
		ChartSubcomponentValue:                         "ClusterIP",
	}
	s.Ports = NewServicePorts()
	return s
}
