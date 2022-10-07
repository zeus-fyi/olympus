package workloads

import (
	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/common"
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/networking"
)

type StatefulSet struct {
	KindDefinition        autogen_structs.ChartComponentKinds
	ParentClassDefinition autogen_structs.ChartSubcomponentParentClassTypes

	Metadata common.Metadata
	Spec     StatefulSetSpec

	//Related Component Requirement
	ServiceDefinition networking.Service
}

type StatefulSetSpec struct {
	Replicas int
	// TODO Selector, VolumeClaimTemplates, ServiceName

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

	return s
}
