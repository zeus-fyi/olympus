package volumes

import (
	v1 "k8s.io/api/core/v1"
)

type VolumeClaimTemplateGroup struct {
	K8sPersistentVolumeClaimSlice []v1.PersistentVolumeClaim
	VolumeClaimTemplateSlice      []VolumeClaimTemplate
}

func NewVolumeClaimTemplateGroup() VolumeClaimTemplateGroup {
	vctg := VolumeClaimTemplateGroup{
		VolumeClaimTemplateSlice: []VolumeClaimTemplate{},
	}
	return vctg
}
