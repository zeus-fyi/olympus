package conversions

import (
	models "github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/charts/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/dev_hacks"
	v1 "k8s.io/api/core/v1"
)

// PodTemplateSpecConfigToDB has a dependency on chart_subcomponent_child_class_types and containers
func PodTemplateSpecConfigToDB(ps *v1.PodTemplateSpec) error {

	zeusTemplateSpec := models.ChartSubcomponentSpecPodTemplateContainers{
		ChartSubcomponentChildClassTypeID: 0, // supplied by chart_subcomponent_child_class_types
		ContainerID:                       0, // supplied by containers
		IsInitContainer:                   false,
	}

	err := dev_hacks.Use(zeusTemplateSpec)
	return err
}
