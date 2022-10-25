package statefulset

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/containers"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/networking"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/common"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type StatefulSetAndChildServices struct {
	StatefulSet StatefulSet
	//Related Component Requirement
	ServiceDefinition networking.Service
}

type StatefulSet struct {
	KindDefinition autogen_bases.ChartComponentResources

	Metadata common.Metadata
	Spec     Spec
}

type Spec struct {
	common.SpecWorkload
	// TODO VolumeClaimTemplates, ServiceName
	Template containers.PodTemplateSpec
}

func (ss *Spec) GetReplicaCount32IntPtr() *int32 {
	return string_utils.ConvertStringTo32BitPtrInt(ss.Replicas.ChartSubcomponentValue)
}

func NewStatefulSet() StatefulSet {
	s := StatefulSet{}
	s.KindDefinition = autogen_bases.ChartComponentResources{
		ChartComponentKindName:   "StatefulSet",
		ChartComponentApiVersion: "apps/v1",
	}
	s.Spec = NewStatefulSetSpec()
	s.Spec.ChartSubcomponentParentClassTypes = autogen_bases.ChartSubcomponentParentClassTypes{
		ChartPackageID:                       0,
		ChartComponentResourceID:             0,
		ChartSubcomponentParentClassTypeID:   0,
		ChartSubcomponentParentClassTypeName: "statefulSetSpec",
	}
	return s
}

func NewStatefulSetSpec() Spec {
	ss := Spec{}
	ss.SpecWorkload = common.NewSpecWorkload()
	ss.Template = containers.NewPodTemplateSpec()
	return ss
}
