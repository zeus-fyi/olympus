package networking

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Service struct {
	K8sService            v1.Service
	KindDefinition        autogen_bases.ChartComponentResources
	ParentClassDefinition autogen_bases.ChartSubcomponentParentClassTypes

	Metadata structs.Metadata
	ServiceSpec
}

func NewService() Service {
	s := Service{}
	typeMeta := metav1.TypeMeta{
		Kind:       "Service",
		APIVersion: "v1",
	}

	s.K8sService = v1.Service{
		TypeMeta:   typeMeta,
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       v1.ServiceSpec{},
		Status:     v1.ServiceStatus{},
	}
	s.KindDefinition = autogen_bases.ChartComponentResources{
		ChartComponentKindName:   "Service",
		ChartComponentApiVersion: "v1",
	}
	s.ServiceSpec = NewServiceSpec()
	return s
}
