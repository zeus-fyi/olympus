package workloads

import (
	autogen_structs2 "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
	common2 "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/common"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/containers"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/networking"
)

type StatefulSetAndChildServices struct {
	StatefulSet StatefulSet
	//Related Component Requirement
	ServiceDefinition networking.Service
}

type StatefulSet struct {
	KindDefinition        autogen_structs2.ChartComponentKinds
	ParentClassDefinition autogen_structs2.ChartSubcomponentParentClassTypes

	Metadata common2.Metadata
	Spec     StatefulSetSpec
}

type StatefulSetSpec struct {
	Replicas int
	Selector common2.Selector
	// TODO VolumeClaimTemplates, ServiceName

	Template containers.PodTemplateSpec
}

func NewStatefulSet() StatefulSet {
	s := StatefulSet{}
	s.KindDefinition = autogen_structs2.ChartComponentKinds{
		ChartComponentKindName:   "StatefulSet",
		ChartComponentApiVersion: "apps/v1",
	}
	s.ParentClassDefinition = autogen_structs2.ChartSubcomponentParentClassTypes{
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
	ss.Selector = common2.NewSelector()
	ss.Template = containers.NewPodTemplateSpec()
	return ss
}
