package ingress

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Ingress struct {
	K8sIngress     v1.Ingress
	KindDefinition autogen_bases.ChartComponentResources

	Metadata structs.ParentMetaData
	Spec
}

func NewIngress() Ingress {
	ing := Ingress{}
	typeMeta := metav1.TypeMeta{
		Kind:       "Ingress",
		APIVersion: "networking.k8s.io/v1",
	}
	ing.K8sIngress = v1.Ingress{
		TypeMeta:   typeMeta,
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       v1.IngressSpec{},
		Status:     v1.IngressStatus{},
	}

	return ing
}
