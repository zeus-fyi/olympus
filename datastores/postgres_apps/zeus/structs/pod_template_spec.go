package structs

import (
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/autogen"
)

type PodTemplateSpec struct {
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

	fields := ChildValuesSlice{}
	genericFields := make(map[string]ChildValuesSlice)
	// (volumes, nodeSelector, affinity, tolerations, etc)
	genericFields["tbd"] = fields

	pts := PodTemplateSpec{
		PodTemplateSpecClassDefinition:    cd,
		PodTemplateSpecClassGenericFields: nil,
	}
	return pts
}
