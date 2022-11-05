package volumes

import "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"

func (v *VolumeClaimTemplateGroup) ConvertK8VolumeClaimTemplateSliceToDB() error {
	v.VolumeClaimTemplateSlice = make([]VolumeClaimTemplate, len(v.K8sPersistentVolumeClaimSlice))
	for i, pvc := range v.K8sPersistentVolumeClaimSlice {
		nPVCDB := NewVolumeClaimTemplate()
		nPVCDB.K8sPersistentVolumeClaim = pvc
		err := nPVCDB.ConvertK8VolumeClaimTemplateSpecToDB()
		if err != nil {
			return err
		}
		nPVCDB.Metadata = common_conversions.ConvertMetadata(pvc.ObjectMeta)
		v.VolumeClaimTemplateSlice[i] = nPVCDB
	}
	return nil
}
