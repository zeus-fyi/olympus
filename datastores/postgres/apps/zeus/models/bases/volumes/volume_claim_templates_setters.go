package volumes

import "github.com/zeus-fyi/olympus/pkg/utils/chronos"

func (v *VolumeClaimTemplateGroup) SetParentIDs(id int) {
	v.ParentClass.SetParentClassTypeID(id)
	for i, _ := range v.VolumeClaimTemplateSlice {
		v.VolumeClaimTemplateSlice[i].Metadata.SetMetadataParentClassTypeIDs(id)
		v.VolumeClaimTemplateSlice[i].Spec.SetParentClassTypeID(id)
	}
}

func (v *VolumeClaimTemplateGroup) SetNewChildIDs() {
	for _, pvc := range v.VolumeClaimTemplateSlice {
		pvc.Spec.SetNewChildClassTypeIDs()
		ts := chronos.Chronos{}
		pvc.Metadata.SetMetadataParentClassTypeIDs(ts.UnixTimeStampNow())
	}
}
