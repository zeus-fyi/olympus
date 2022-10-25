package containers

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	cont_conv "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/containers"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/common"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/containers"
	v1 "k8s.io/api/core/v1"
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
	PodTemplateContainers             containers.Containers
}

func (p *PodTemplateSpec) GetContainers() containers.Containers {
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

func (p *PodTemplateSpec) AddContainer(c containers.Container) {
	c.ClassDefinition = autogen_bases.ChartSubcomponentChildClassTypes{
		ChartSubcomponentChildClassTypeID:   p.Spec.PodTemplateSpecClassDefinition.ChartSubcomponentChildClassTypeID,
		ChartSubcomponentChildClassTypeName: "container",
	}
	p.Spec.PodTemplateContainers = append(p.Spec.PodTemplateContainers, c)
}

func (p *PodTemplateSpec) GetPodSpecParentClassTypeID() int {
	return p.Spec.PodTemplateSpecClassDefinition.ChartSubcomponentParentClassTypeID
}

func (p *PodTemplateSpec) SetPodSpecParentClassTypeID(id int) {
	p.Spec.PodTemplateSpecClassDefinition.ChartSubcomponentParentClassTypeID = id
}

func (p *PodTemplateSpec) GetPodSpecChildClassTypeID() int {
	return p.Spec.PodTemplateSpecClassDefinition.ChartSubcomponentChildClassTypeID
}
func (p *PodTemplateSpec) SetPodSpecChildClassTypeID(id int) {
	p.Spec.PodTemplateSpecClassDefinition.ChartSubcomponentChildClassTypeID = id
}

// ConvertPodTemplateSpecConfigToDB PodTemplateSpecConfigToDB has a dependency on chart_subcomponent_child_class_types and containers
func (p *PodTemplateSpec) ConvertPodTemplateSpecConfigToDB(ps *v1.PodSpec) (PodTemplateSpec, error) {
	dbPodSpec := NewPodTemplateSpec()

	dbSpecVolumes, err := common_conversions.VolumesToDB(ps.Volumes)
	if err != nil {
		return dbPodSpec, err
	}
	dbSpecContainers, err := cont_conv.ConvertContainersToDB(ps.Containers)
	if err != nil {
		return dbPodSpec, err
	}
	dbPodSpec.Spec.PodTemplateContainers = dbSpecContainers
	dbPodSpec.Spec.PodTemplateSpecVolumes = dbSpecVolumes
	if err != nil {
		return dbPodSpec, err
	}
	return dbPodSpec, nil
}
