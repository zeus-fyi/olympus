package volumes

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	v1 "k8s.io/api/core/v1"
)

type VolumeClaimTemplate struct {
	K8sPersistentVolumeClaim v1.PersistentVolumeClaim
	Spec                     VolumeClaimTemplateSpec
	Metadata                 structs.Metadata
}

func NewVolumeClaimTemplate() VolumeClaimTemplate {
	vct := VolumeClaimTemplate{
		Spec: NewVolumeClaimTemplateSpec(),
	}
	return vct
}
