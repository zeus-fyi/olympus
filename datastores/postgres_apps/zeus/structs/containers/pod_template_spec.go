package containers

import (
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/common"
)

type PodTemplateSpec struct {
	Metadata common.Metadata
	Spec     PodSpec
}

type PodSpec struct {
	PodTemplateSpecClassDefinition    autogen_structs.ChartSubcomponentChildClassTypes
	PodTemplateSpecClassGenericFields map[string]common.ChildValuesSlice
	PodTemplateSpecVolumes            common.VolumesSlice
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
		Metadata: common.Metadata{},
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
