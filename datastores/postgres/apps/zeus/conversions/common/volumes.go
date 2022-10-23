package common

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	v1 "k8s.io/api/core/v1"
)

func VolumesToDB(vs []v1.Volume) []autogen_bases.Volumes {
	dbVolSlice := make([]autogen_bases.Volumes, len(vs))
	for i, v := range vs {
		dbVol := VolumeToDB(v)
		dbVolSlice[i] = dbVol
	}
	return dbVolSlice
}

func VolumeToDB(v v1.Volume) autogen_bases.Volumes {
	dbVolume := autogen_bases.Volumes{
		VolumeID:             0,
		VolumeName:           v.Name,
		VolumeKeyValuesJSONb: "",
	}
	return dbVolume
}
