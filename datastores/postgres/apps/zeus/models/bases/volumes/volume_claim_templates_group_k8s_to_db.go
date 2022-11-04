package volumes

func (v *VolumeClaimTemplateGroup) ConvertK8VolumeClaimTemplateSliceToDB() error {
	v.VolumeClaimTemplateSlice = make([]VolumeClaimTemplate, len(v.K8sPersistentVolumeClaimSlice))
	for i, pvc := range v.K8sPersistentVolumeClaimSlice {
		nPVCDB := NewVolumeClaimTemplate()
		nPVCDB.K8sPersistentVolumeClaim = pvc
		err := nPVCDB.ConvertK8VolumeClaimTemplateSpecToDB()
		if err != nil {
			return err
		}
		v.VolumeClaimTemplateSlice[i] = nPVCDB
	}
	return nil
}
