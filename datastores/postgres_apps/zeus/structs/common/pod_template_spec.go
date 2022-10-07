package common

import (
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/autogen"
)

type PodTemplateSpec struct {
	Metadata Metadata
	Spec     PodSpec
}

type PodSpec struct {
	PodTemplateSpecClassDefinition    autogen_structs.ChartSubcomponentChildClassTypes
	PodTemplateSpecClassGenericFields map[string]ChildValuesSlice
	PodTemplateSpecVolumes            VolumesSlice
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
		Metadata: Metadata{},
		Spec:     ps,
	}

	return pts
}
