package volumes

import "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"

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
