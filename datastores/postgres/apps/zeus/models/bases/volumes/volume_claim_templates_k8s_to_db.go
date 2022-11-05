package volumes

import "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"

func (v *VolumeClaimTemplate) ConvertK8VolumeClaimTemplateToDB() error {
	meta := v.K8sPersistentVolumeClaim.ObjectMeta
	v.Metadata.Metadata = common_conversions.CreateMetadataByFields(meta.Name, meta.Annotations, meta.Labels)
	err := v.ConvertK8VolumeClaimTemplateSpecToDB()
	return err
}
