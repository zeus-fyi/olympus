package containers

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/common"
)

type PodTemplateSpec struct {
	autogen_bases.ChartSubcomponentChildClassTypes
	Metadata common.Metadata
	Spec     PodSpec
}

type PodSpec struct {
	PodTemplateSpecClassDefinition    autogen_bases.ChartSubcomponentChildClassTypes
	PodTemplateSpecClassGenericFields map[string]common.ChildValuesSlice
	PodTemplateSpecVolumes            autogen_bases.VolumesSlice
	PodTemplateContainers             Containers
}

func (p *PodTemplateSpec) GetContainers() Containers {
	return p.Spec.PodTemplateContainers
}

func NewPodTemplateSpec() PodTemplateSpec {
	cd := autogen_bases.ChartSubcomponentChildClassTypes{
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
	c.ClassDefinition = autogen_bases.ChartSubcomponentChildClassTypes{
		ChartSubcomponentChildClassTypeID:   p.Spec.PodTemplateSpecClassDefinition.ChartSubcomponentChildClassTypeID,
		ChartSubcomponentChildClassTypeName: "container",
	}
	p.Spec.PodTemplateContainers = append(p.Spec.PodTemplateContainers, c)
}

func (p *PodTemplateSpec) GetPodSpecChildClassTypeID() int {
	return p.Spec.PodTemplateSpecClassDefinition.ChartSubcomponentChildClassTypeID
}
