package workloads

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/common"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/containers"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/networking"
)

type StatefulSetAndChildServices struct {
	StatefulSet StatefulSet
	//Related Component Requirement
	ServiceDefinition networking.Service
}

type StatefulSet struct {
	KindDefinition        autogen_bases.ChartComponentResources
	ParentClassDefinition autogen_bases.ChartSubcomponentParentClassTypes

	Metadata common.Metadata
	Spec     StatefulSetSpec
}

type StatefulSetSpec struct {
	Replicas int
	Selector common.Selector
	// TODO VolumeClaimTemplates, ServiceName

	Template containers.PodTemplateSpec
}

func NewStatefulSet() StatefulSet {
	s := StatefulSet{}
	s.KindDefinition = autogen_bases.ChartComponentResources{
		ChartComponentKindName:   "StatefulSet",
		ChartComponentApiVersion: "apps/v1",
	}
	s.ParentClassDefinition = autogen_bases.ChartSubcomponentParentClassTypes{
		ChartPackageID:                       0,
		ChartComponentResourceID:             0,
		ChartSubcomponentParentClassTypeID:   0,
		ChartSubcomponentParentClassTypeName: "statefulSetSpec",
	}

	s.Spec = NewStatefulSetSpec()
	return s
}

func NewStatefulSetSpec() StatefulSetSpec {
	ss := StatefulSetSpec{}
	ss.Selector = common.NewSelector()
	ss.Template = containers.NewPodTemplateSpec()
	return ss
}
