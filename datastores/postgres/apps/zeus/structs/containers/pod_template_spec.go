package containers

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
	common2 "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/common"
)

type PodTemplateSpec struct {
	Metadata common2.Metadata
	Spec     PodSpec
}

type PodSpec struct {
	PodTemplateSpecClassDefinition    autogen_structs.autogen_structs
	PodTemplateSpecClassGenericFields map[string]common2.ChildValuesSlice
	PodTemplateSpecVolumes            common2.VolumesSlice
	PodTemplateContainers             Containers
}

func NewPodTemplateSpec() PodTemplateSpec {
	cd := autogen_structs.ChartSubcomponentChildClassTypes{
		ChartSubcomponentParentClassTypeID:  0,
		ChartSubcomponentChildClassTypeID:   0,
		ChartSubcomponentChildClassTypeName: "PodTemplateSpec",
	}

	ps := PodSpec{
		PodTemplateSpecClassDefinition:    cd,
		PodTemplateSpecClassGenericFields: nil,
	}

	pts := PodTemplateSpec{
		Metadata: common2.Metadata{},
		Spec:     ps,
	}

	return pts
}

func (p *PodTemplateSpec) AddContainer(c Container) {
	c.ClassDefinition = autogen_structs.ChartSubcomponentChildClassTypes{
		ChartSubcomponentChildClassTypeID:   p.Spec.PodTemplateSpecClassDefinition.ChartSubcomponentChildClassTypeID,
		ChartSubcomponentChildClassTypeName: "container",
	}
	p.Spec.PodTemplateContainers = append(p.Spec.PodTemplateContainers, c)
}
