package containers

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	cont_conv "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/containers"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	v1 "k8s.io/api/core/v1"
)

// ConvertPodTemplateSpecConfigToDB PodTemplateSpecConfigToDB has a dependency on chart_subcomponent_child_class_types and containers
func (p *PodTemplateSpec) ConvertPodTemplateSpecConfigToDB(ps *v1.PodSpec) error {
	dbSpecVolumes, err := common_conversions.VolumesToDB(ps.Volumes)
	if err != nil {
		return err
	}
	dbSpecContainers, err := cont_conv.ConvertContainersToDB(ps.Containers, false)
	if err != nil {
		return err
	}
	dbSpecInitContainers, err := cont_conv.ConvertContainersToDB(ps.InitContainers, true)
	if err != nil {
		return err
	}

	if ps.TerminationGracePeriodSeconds != nil {
		gps := string_utils.Convert64BitPtrIntToString(ps.TerminationGracePeriodSeconds)
		csv := structs.ChildClassSingleValue{
			ChartSubcomponentChildClassTypes: autogen_bases.ChartSubcomponentChildClassTypes{},
			ChartSubcomponentsChildValues: autogen_bases.ChartSubcomponentsChildValues{
				ChartSubcomponentChildClassTypeID:              0,
				ChartSubcomponentChartPackageTemplateInjection: false,
				ChartSubcomponentKeyName:                       "terminationGracePeriodSeconds",
				ChartSubcomponentValue:                         gps,
			},
		}
		p.AddPodTemplateSpecClassGenericFields(csv)
	}

	dbSpecInitContainers = append(dbSpecInitContainers, dbSpecContainers...)
	p.Spec.PodTemplateContainers = dbSpecInitContainers
	p.Spec.PodTemplateSpecVolumes = dbSpecVolumes
	if err != nil {
		return err
	}

	return nil
}
