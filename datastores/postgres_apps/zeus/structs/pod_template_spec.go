package structs

import (
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/common"
)

type PodTemplateSpec struct {
	PodTemplateSpecClassDefinition    autogen_structs.ChartSubcomponentChildClassTypes
	PodTemplateSpecClassGenericFields map[string]common.ChildValuesSlice
	PodTemplateSpecVolumes            common.VolumesSlice
	PodTemplateContainers             common.Containers
}

func NewPodTemplateSpec() PodTemplateSpec {
	cd := autogen_structs.ChartSubcomponentChildClassTypes{
		ChartSubcomponentParentClassTypeID:  0,
		ChartSubcomponentChildClassTypeID:   0,
		ChartSubcomponentChildClassTypeName: "PodTemplateSpec",
	}

	fields := common.ChildValuesSlice{}
	genericFields := make(map[string]common.ChildValuesSlice)
	// (volumes, nodeSelector, affinity, tolerations, etc)
	genericFields["tbd"] = fields

	pts := PodTemplateSpec{
		PodTemplateSpecClassDefinition:    cd,
		PodTemplateSpecClassGenericFields: nil,
	}
	return pts
}
