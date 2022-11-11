package volumes

import "encoding/json"

func (v *VolumeClaimTemplate) ConvertK8VolumeClaimTemplateSpecToDB() error {
	v.ConvertK8VolumeClaimTemplateSpecStorageClassNameToDB()
	v.ConvertK8VolumeClaimTemplateSpecAccessModesToDB()
	err := v.ConvertK8VolumeClaimTemplateSpecResourceRequestsToDB()
	return err
}

func (v *VolumeClaimTemplate) ConvertK8VolumeClaimTemplateSpecStorageClassNameToDB() {
	v.Spec.StorageClassName.ChartSubcomponentChildClassTypeName = "VolumeClaimTemplateSpec"
	scName := v.K8sPersistentVolumeClaim.Spec.StorageClassName
	if scName != nil {
		v.Spec.StorageClassName.ChartSubcomponentKeyName = "storageClassName"
		v.Spec.StorageClassName.ChartSubcomponentValue = *scName
	}
}

func (v *VolumeClaimTemplate) ConvertK8VolumeClaimTemplateSpecAccessModesToDB() {
	v.Spec.AccessModes.ChartSubcomponentChildClassTypeName = "VolumeClaimTemplateSpec"
	accessModes := v.K8sPersistentVolumeClaim.Spec.AccessModes
	for _, am := range accessModes {
		v.Spec.AccessModes.AddKeyValue("accessMode", string(am))
	}
}

func (v *VolumeClaimTemplate) ConvertK8VolumeClaimTemplateSpecResourceRequestsToDB() error {
	v.Spec.ResourceRequests.ChartSubcomponentChildClassTypeName = "VolumeClaimTemplateSpec"
	rr := v.K8sPersistentVolumeClaim.Spec.Resources
	for _, r := range rr.Limits {
		b, err := json.Marshal(r)
		if err != nil {
			return err
		}
		v.Spec.ResourceRequests.AddKeyValue("limits", string(b))
	}
	for _, r := range rr.Requests {
		b, err := json.Marshal(r)
		if err != nil {
			return err
		}
		v.Spec.ResourceRequests.AddKeyValue("requests", string(b))
	}
	return nil
}
