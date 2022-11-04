package volumes

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	v1 "k8s.io/api/core/v1"
)

const PVCChartComponentResourceID = 2

type VolumeClaimTemplate struct {
	K8sPersistentVolumeClaim v1.PersistentVolumeClaim
	Spec                     VolumeClaimTemplateSpec
	Metadata                 structs.Metadata
}

type VolumeClaimTemplateGroup struct {
	common.ParentClass
	K8sPersistentVolumeClaimSlice []v1.PersistentVolumeClaim
	VolumeClaimTemplateSlice      []VolumeClaimTemplate
}

func NewVolumeClaimTemplateGroup() VolumeClaimTemplateGroup {
	pc := common.ParentClass{ChartSubcomponentParentClassTypes: autogen_bases.ChartSubcomponentParentClassTypes{
		ChartPackageID:                       0,
		ChartComponentResourceID:             PVCChartComponentResourceID,
		ChartSubcomponentParentClassTypeID:   0,
		ChartSubcomponentParentClassTypeName: "VolumeClaimTemplateSpec",
	}}

	vctg := VolumeClaimTemplateGroup{
		ParentClass:              pc,
		VolumeClaimTemplateSlice: []VolumeClaimTemplate{},
	}
	return vctg
}

func NewVolumeClaimTemplate() VolumeClaimTemplate {
	vct := VolumeClaimTemplate{
		Spec: NewVolumeClaimTemplateSpec(),
	}
	return vct
}
