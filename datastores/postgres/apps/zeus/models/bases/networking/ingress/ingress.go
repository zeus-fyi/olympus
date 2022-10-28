package ingress

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	v1 "k8s.io/api/networking/v1"
)

const IngressChartComponentResourceID = 14

type Ingress struct {
	K8sIngress     v1.Ingress
	KindDefinition autogen_bases.ChartComponentResources

	Metadata structs.ParentMetaData
	Spec
}

func NewIngress() Ingress {
	ing := Ingress{}
	ing.K8sIngress = v1.Ingress{}
	ing.KindDefinition = autogen_bases.ChartComponentResources{
		ChartComponentKindName:   "Ingress",
		ChartComponentApiVersion: "networking.k8s.io/v1",
		ChartComponentResourceID: IngressChartComponentResourceID,
	}
	ing.Metadata.ChartComponentResourceID = IngressChartComponentResourceID
	ing.Metadata.ChartSubcomponentParentClassTypeName = "IngressParentMetadata"
	ing.Metadata.Metadata = structs.NewMetadata()
	ing.Spec = NewIngressSpec()
	return ing
}
