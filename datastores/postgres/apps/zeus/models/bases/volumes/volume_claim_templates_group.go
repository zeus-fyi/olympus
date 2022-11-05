package volumes

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	v1 "k8s.io/api/core/v1"
)

type VolumeClaimTemplateGroup struct {
	common.ParentClass
	K8sPersistentVolumeClaimSlice []v1.PersistentVolumeClaim
	VolumeClaimTemplateSlice      []VolumeClaimTemplate
}

func NewVolumeClaimTemplateGroup() VolumeClaimTemplateGroup {
	pc := common.ParentClass{ChartSubcomponentParentClassTypes: autogen_bases.ChartSubcomponentParentClassTypes{
		ChartPackageID:                       0,
		ChartSubcomponentParentClassTypeID:   0,
		ChartSubcomponentParentClassTypeName: "VolumeClaimTemplateSpec",
	}}

	vctg := VolumeClaimTemplateGroup{
		ParentClass:              pc,
		VolumeClaimTemplateSlice: []VolumeClaimTemplate{},
	}
	return vctg
}
