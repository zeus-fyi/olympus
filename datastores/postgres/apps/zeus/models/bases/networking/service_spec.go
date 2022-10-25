package networking

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs/common"
)

// ServiceSpec has these type options: ClusterIP, NodePort, LoadBalancer, ExternalName
type ServiceSpec struct {
	Type     autogen_bases.ChartSubcomponentsChildValues
	Selector common.Selector
	Ports    ServicePorts
}

func NewServiceSpec() ServiceSpec {
	s := ServiceSpec{}
	s.Type = autogen_bases.ChartSubcomponentsChildValues{
		ChartSubcomponentChildClassTypeID:              0,
		ChartSubcomponentChartPackageTemplateInjection: true,
		ChartSubcomponentKeyName:                       "type",
		ChartSubcomponentValue:                         "ClusterIP",
	}
	s.Ports = NewServicePorts()
	return s
}
