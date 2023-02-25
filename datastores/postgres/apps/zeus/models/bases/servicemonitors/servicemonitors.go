package servicemonitors

import (
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const ServiceMonitorChartComponentResourceID = 27

type ServiceMonitor struct {
	K8sServiceMonitor v1.ServiceMonitor
	KindDefinition    autogen_bases.ChartComponentResources

	Metadata structs.ParentMetaData
	Spec     ServiceMonitorSpec
}

type ServiceMonitorSpec struct {
	common.ParentClass
	structs.ChildClassSingleValue
}

func NewServiceMonitor() ServiceMonitor {
	s := ServiceMonitor{}
	typeMeta := metav1.TypeMeta{
		Kind:       "ServiceMonitor",
		APIVersion: "monitoring.coreos.com/v1",
	}
	s.K8sServiceMonitor = v1.ServiceMonitor{TypeMeta: typeMeta}
	s.KindDefinition = autogen_bases.ChartComponentResources{
		ChartComponentKindName:   "ServiceMonitor",
		ChartComponentApiVersion: "monitoring.coreos.com/v1",
		ChartComponentResourceID: ServiceMonitorChartComponentResourceID,
	}
	s.Spec.ChartSubcomponentParentClassTypes = autogen_bases.ChartSubcomponentParentClassTypes{
		ChartPackageID:                       0,
		ChartComponentResourceID:             ServiceMonitorChartComponentResourceID,
		ChartSubcomponentParentClassTypeID:   0,
		ChartSubcomponentParentClassTypeName: "Spec",
	}
	s.Metadata.Metadata = structs.NewMetadata()
	s.Metadata.ChartSubcomponentParentClassTypeName = "ServiceMonitorParentMetadata"
	s.Metadata.ChartComponentResourceID = ServiceMonitorChartComponentResourceID
	return s
}
