package volumes

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/common"
	v1 "k8s.io/api/core/v1"
)

type VolumeClaimTemplate struct {
	K8sPersistentVolumeClaim v1.PersistentVolumeClaim
	common.ParentClass
	Spec     VolumeClaimTemplateSpec
	Metadata structs.ParentMetaData
}

func NewVolumeClaimTemplate() VolumeClaimTemplate {
	pc := common.ParentClass{ChartSubcomponentParentClassTypes: autogen_bases.ChartSubcomponentParentClassTypes{
		ChartPackageID:                       0,
		ChartSubcomponentParentClassTypeID:   0,
		ChartSubcomponentParentClassTypeName: "VolumeClaimTemplate",
	}}
	vct := VolumeClaimTemplate{
		ParentClass: pc,
		Spec:        NewVolumeClaimTemplateSpec(),
	}
	vct.Metadata.Metadata = structs.NewMetadata()
	vct.Metadata.ChartSubcomponentParentClassTypeName = "VolumeClaimTemplate"
	return vct
}
