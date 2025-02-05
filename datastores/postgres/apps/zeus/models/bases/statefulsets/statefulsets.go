package statefulsets

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking/services"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/volumes"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/containers"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const StsChartComponentResourceID = 1

type StatefulSetAndChildServices struct {
	StatefulSet StatefulSet
	//Related Component Requirement
	ServiceDefinition services.Service
}

type StatefulSet struct {
	K8sStatefulSet v1.StatefulSet
	KindDefinition autogen_bases.ChartComponentResources

	Metadata structs.ParentMetaData
	Spec     Spec
}

type Spec struct {
	structs.SpecWorkload
	Template containers.PodTemplateSpec

	StatefulSetUpdateStrategy structs.ChildClassMultiValue
	PodManagementPolicy       structs.ChildClassSingleValue
	ServiceName               structs.ChildClassSingleValue

	VolumeClaimTemplates volumes.VolumeClaimTemplateGroup
}

func NewStatefulSet() StatefulSet {
	s := StatefulSet{}
	typeMeta := metav1.TypeMeta{
		Kind:       "StatefulSet",
		APIVersion: "apps/v1",
	}
	s.K8sStatefulSet = v1.StatefulSet{TypeMeta: typeMeta}
	s.KindDefinition = autogen_bases.ChartComponentResources{
		ChartComponentKindName:   "StatefulSet",
		ChartComponentApiVersion: "apps/v1",
		ChartComponentResourceID: StsChartComponentResourceID,
	}
	s.Spec = NewStatefulSetSpec()
	s.Spec.ChartSubcomponentParentClassTypes = autogen_bases.ChartSubcomponentParentClassTypes{
		ChartPackageID:                       0,
		ChartComponentResourceID:             StsChartComponentResourceID,
		ChartSubcomponentParentClassTypeID:   0,
		ChartSubcomponentParentClassTypeName: "Spec",
	}
	s.Metadata.Metadata = structs.NewMetadata()
	s.Metadata.ChartSubcomponentParentClassTypeName = "StatefulSetSpecParentMetadata"
	s.Metadata.ChartComponentResourceID = StsChartComponentResourceID
	s.Spec.VolumeClaimTemplates = volumes.NewVolumeClaimTemplateGroup()
	return s
}

func NewStatefulSetSpec() Spec {
	ss := Spec{}
	ss.SpecWorkload = structs.NewSpecWorkload()
	ss.Template = containers.NewPodTemplateSpec()
	ss.SpecWorkload.ChartSubcomponentParentClassTypeName = "Spec"
	return ss
}
