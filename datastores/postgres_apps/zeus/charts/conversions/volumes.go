package conversions

import (
	v1 "k8s.io/api/core/v1"

	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/charts/structs/autogen"
)

func VolumesToDB(vs []v1.Volume) []autogen_structs.Volumes {
	dbVolSlice := make([]autogen_structs.Volumes, len(vs))
	for i, v := range vs {
		dbVol := VolumeToDB(v)
		dbVolSlice[i] = dbVol
	}
	return dbVolSlice
}

func VolumeToDB(v v1.Volume) autogen_structs.Volumes {
	dbVolume := autogen_structs.Volumes{
		VolumeID:             0,
		VolumeName:           v.Name,
		VolumeKeyValuesJSONb: "",
	}
	return dbVolume
}
