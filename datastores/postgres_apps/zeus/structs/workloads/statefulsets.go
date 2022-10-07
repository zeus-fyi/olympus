package workloads

import (
	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/common"
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/networking"
)

type StatefulSetAndChildServices struct {
	StatefulSet StatefulSet
	//Related Component Requirement
	ServiceDefinition networking.Service
}

type StatefulSet struct {
	KindDefinition        autogen_structs.ChartComponentKinds
	ParentClassDefinition autogen_structs.ChartSubcomponentParentClassTypes

	Metadata common.Metadata
	Spec     StatefulSetSpec
}

type StatefulSetSpec struct {
	Replicas int
	Selector common.Selector
	// TODO VolumeClaimTemplates, ServiceName

	Template common.PodTemplateSpec
}

func NewStatefulSet() StatefulSet {
	s := StatefulSet{}
	s.KindDefinition = autogen_structs.ChartComponentKinds{
		ChartComponentKindName:   "StatefulSet",
		ChartComponentApiVersion: "apps/v1",
	}
	s.ParentClassDefinition = autogen_structs.ChartSubcomponentParentClassTypes{
		ChartPackageID:                       0,
		ChartComponentKindID:                 0,
		ChartSubcomponentParentClassTypeID:   0,
		ChartSubcomponentParentClassTypeName: "statefulSetSpec",
	}

	s.Spec = NewStatefulSetSpec()
	return s
}

func NewStatefulSetSpec() StatefulSetSpec {
	ss := StatefulSetSpec{}
	ss.Selector = common.NewSelector()
	ss.Template = common.NewPodTemplateSpec()
	return ss
}
