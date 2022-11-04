package volumes

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

type VolumeClaimTemplateSpec struct {
	StorageClassName structs.ChildClassSingleValue
	AccessModes      structs.ChildClassMultiValue
	ResourceRequests structs.ChildClassMultiValue
}

func NewVolumeClaimTemplateSpec() VolumeClaimTemplateSpec {
	vtcs := VolumeClaimTemplateSpec{
		StorageClassName: structs.ChildClassSingleValue{},
		AccessModes:      structs.ChildClassMultiValue{},
		ResourceRequests: structs.ChildClassMultiValue{},
	}
	return vtcs
}

func (v *VolumeClaimTemplateSpec) SetParentClassTypeID(id int) {
	v.StorageClassName.SetParentClassTypeID(id)
	v.AccessModes.SetParentClassTypeID(id)
	v.ResourceRequests.SetParentClassTypeID(id)
}

func (v *VolumeClaimTemplateSpec) SetNewChildClassTypeIDs() {
	ts := chronos.Chronos{}
	v.StorageClassName.SetChildClassTypeIDs(ts.UnixTimeStampNow())
	v.AccessModes.SetChildClassTypeIDs(ts.UnixTimeStampNow())
	v.ResourceRequests.SetChildClassTypeIDs(ts.UnixTimeStampNow())
}
