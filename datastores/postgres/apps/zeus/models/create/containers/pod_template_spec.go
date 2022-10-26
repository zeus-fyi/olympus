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
	PodTemplateMapK8sContainers       map[int]v1.Container
}

func (p *PodTemplateSpec) GetContainerMap(id int) v1.Container {
	if _, ok := p.Spec.PodTemplateMapK8sContainers[id]; !ok {
		return v1.Container{}
	}
	return p.Spec.PodTemplateMapK8sContainers[id]
}

func (p *PodTemplateSpec) SetContainerMap(id int, c v1.Container) {
	p.Spec.PodTemplateMapK8sContainers[id] = c
}

func (p *PodTemplateSpec) AddVolume(v autogen_bases.Volumes) {
	p.Spec.PodTemplateSpecVolumes = append(p.Spec.PodTemplateSpecVolumes, v)
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
		PodTemplateMapK8sContainers:       make(map[int]v1.Container),
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
