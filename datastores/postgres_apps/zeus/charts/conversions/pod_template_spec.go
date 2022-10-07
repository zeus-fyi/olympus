package conversions

import (
	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/charts/structs/autogen"
	v1 "k8s.io/api/core/v1"

	"github.com/zeus-fyi/olympus/pkg/utils/dev_hacks"
)

// PodTemplateSpecConfigToDB has a dependency on chart_subcomponent_child_class_types and containers
func PodTemplateSpecConfigToDB(ps *v1.PodTemplateSpec) error {

	cp := autogen_structs.ChartPackages{
		ChartPackageID: 0,
		ChartName:      "",
		ChartVersion:   "",
	}

	cpr := autogen_structs.ChartPackageComponents{
		ChartSubcomponentParentClassTypeID: 0,
	}
	cpk := autogen_structs.ChartComponentKinds{
		ChartComponentKindID:     0,
		ChartComponentKindName:   "",
		ChartComponentApiVersion: "",
	}

	pc := autogen_structs.ChartSubcomponentParentClassTypes{
		ChartPackageID:                       0,
		ChartComponentKindID:                 0,
		ChartSubcomponentParentClassTypeID:   0,
		ChartSubcomponentParentClassTypeName: "",
	}

	cct := autogen_structs.ChartSubcomponentChildClassTypes{
		ChartSubcomponentParentClassTypeID:  0,
		ChartSubcomponentChildClassTypeID:   0,
		ChartSubcomponentChildClassTypeName: "",
	}

	zeusTemplateSpec := autogen_structs.ChartSubcomponentSpecPodTemplateContainers{
		ChartSubcomponentChildClassTypeID: 0,
		ContainerID:                       0,
		IsInitContainer:                   false,
		ContainerSortOrder:                0,
	}
	_ = dev_hacks.Use(cpk)

	_ = dev_hacks.Use(cpr)

	_ = dev_hacks.Use(cp)

	_ = dev_hacks.Use(cct)
	err := dev_hacks.Use(pc)
	err = dev_hacks.Use(zeusTemplateSpec)

	return err
}
