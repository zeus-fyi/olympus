package containers

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/containers"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	v1 "k8s.io/api/core/v1"
)

type PodTemplateSpec struct {
	common.ParentClass
	Metadata structs.ParentMetaData
	Spec     PodSpec
}

type PodSpec struct {
	K8sPodSpec *v1.PodSpec

	PodTemplateSpecClassDefinition    autogen_bases.ChartSubcomponentChildClassTypes
	PodTemplateSpecClassGenericFields map[string]structs.ChildClassSingleValue
	PodTemplateSpecVolumes            autogen_bases.VolumesSlice
	PodTemplateContainers             containers.Containers

	ShareProcessNamespace *structs.ChildClassSingleValue
}

func (p *PodTemplateSpec) AddPodTemplateSpecClassGenericFields(cv structs.ChildClassSingleValue) {
	cv.ChartSubcomponentChildClassTypeName = "PodTemplateSpec"
	p.Spec.PodTemplateSpecClassGenericFields[cv.ChartSubcomponentKeyName] = cv
}

func (p *PodTemplateSpec) SetK8sPodSpecVolumes(vs []v1.Volume) {
	p.Spec.K8sPodSpec.Volumes = vs
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
		PodTemplateSpecClassGenericFields: make(map[string]structs.ChildClassSingleValue),
		K8sPodSpec:                        &v1.PodSpec{},
	}

	pts := PodTemplateSpec{
		ParentClass: common.ParentClass{
			ChartSubcomponentParentClassTypes: autogen_bases.ChartSubcomponentParentClassTypes{
				ChartPackageID:                       0,
				ChartComponentResourceID:             0,
				ChartSubcomponentParentClassTypeID:   0,
				ChartSubcomponentParentClassTypeName: "PodTemplateSpec",
			}},
		Metadata: structs.NewParentMetaData("PodTemplateSpecMetadata"),
		Spec:     ps,
	}
	return pts
}

func (p *PodTemplateSpec) GetPodSpecParentClassTypeID() int {
	return p.ChartSubcomponentParentClassTypeID
}

func (p *PodTemplateSpec) SetPodSpecParentClassTypeID(id int) {
	p.ParentClass.SetParentClassTypeID(id)
	p.Spec.PodTemplateSpecClassDefinition.ChartSubcomponentParentClassTypeID = id
}

func (p *PodTemplateSpec) GetPodSpecChildClassTypeID() int {
	return p.Spec.PodTemplateSpecClassDefinition.ChartSubcomponentChildClassTypeID
}

func (p *PodTemplateSpec) SetPodSpecChildClassTypeID(id int) {
	p.Spec.PodTemplateSpecClassDefinition.ChartSubcomponentChildClassTypeID = id
}
