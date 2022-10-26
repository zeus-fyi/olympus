package containers

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	cont_conv "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/containers"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/containers"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	v1 "k8s.io/api/core/v1"
)

type PodTemplateSpec struct {
	common.ParentClass
	Metadata structs.Metadata
	Spec     PodSpec
}

type PodSpec struct {
	PodTemplateSpecClassDefinition    autogen_bases.ChartSubcomponentChildClassTypes
	PodTemplateSpecClassGenericFields map[string]structs.ChildValuesSlice
	PodTemplateSpecVolumes            autogen_bases.VolumesSlice
	PodTemplateContainers             containers.Containers
}

func (p *PodTemplateSpec) GetContainers() containers.Containers {
	return p.Spec.PodTemplateContainers
}

func (p *PodTemplateSpec) AddContainer(c containers.Container) {
	p.Spec.PodTemplateContainers = append(p.Spec.PodTemplateContainers, c)
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
		ParentClass: common.ParentClass{
			ChartSubcomponentParentClassTypes: autogen_bases.ChartSubcomponentParentClassTypes{
				ChartPackageID:                       0,
				ChartComponentResourceID:             0,
				ChartSubcomponentParentClassTypeID:   0,
				ChartSubcomponentParentClassTypeName: "PodTemplateSpec",
			}},
		Metadata: structs.Metadata{},
		Spec:     ps,
	}

	return pts
}

func (p *PodTemplateSpec) GetPodSpecParentClassTypeID() int {
	return p.ChartSubcomponentParentClassTypeID
}

func (p *PodTemplateSpec) SetPodSpecParentClassTypeID(id int) {
	p.ParentClass.InsertParentClassTypeID(id)
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
