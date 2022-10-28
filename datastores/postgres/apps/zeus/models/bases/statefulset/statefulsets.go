package statefulset

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking/services"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/containers"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StatefulSetAndChildServices struct {
	StatefulSet StatefulSet
	//Related Component Requirement
	ServiceDefinition services.Service
}

type StatefulSet struct {
	K8sDeployment *v1.StatefulSet

	KindDefinition autogen_bases.ChartComponentResources

	Metadata structs.ParentMetaData
	Spec     Spec
}

type Spec struct {
	structs.SpecWorkload
	// TODO VolumeClaimTemplates, ServiceName
	Template containers.PodTemplateSpec
}

func NewStatefulSet() StatefulSet {
	s := StatefulSet{}
	typeMeta := metav1.TypeMeta{
		Kind:       "StatefulSet",
		APIVersion: "apps/v1",
	}
	s.K8sDeployment = &v1.StatefulSet{TypeMeta: typeMeta}
	s.KindDefinition = autogen_bases.ChartComponentResources{
		ChartComponentKindName:   "StatefulSet",
		ChartComponentApiVersion: "apps/v1",
	}
	s.Spec = NewStatefulSetSpec()
	s.Spec.ChartSubcomponentParentClassTypes = autogen_bases.ChartSubcomponentParentClassTypes{
		ChartPackageID:                       0,
		ChartComponentResourceID:             0,
		ChartSubcomponentParentClassTypeID:   0,
		ChartSubcomponentParentClassTypeName: "StatefulSetSpec",
	}
	s.Metadata.Metadata = structs.NewMetadata()
	s.Metadata.ChartSubcomponentParentClassTypeName = "StatefulSetSpecParentMetadata"
	return s
}

func NewStatefulSetSpec() Spec {
	ss := Spec{}
	ss.SpecWorkload = structs.NewSpecWorkload()
	ss.Template = containers.NewPodTemplateSpec()
	ss.SpecWorkload.ChartSubcomponentParentClassTypeName = "Spec"
	return ss
}
