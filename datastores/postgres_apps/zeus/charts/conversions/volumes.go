package conversions

import (
	v1 "k8s.io/api/core/v1"

	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/charts/structs/autogen"
)

func VolumeToDB(v *v1.Volume) autogen_structs.Volumes {
	dbContainer := autogen_structs.Volumes{
		VolumeID:             0,
		VolumeName:           v.Name,
		VolumeKeyValuesJSONb: "",
	}
	return dbContainer
}
