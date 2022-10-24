package common_conversions

import (
	"encoding/json"

	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	v1 "k8s.io/api/core/v1"
)

func VolumesToDB(vs []v1.Volume) ([]autogen_bases.Volumes, error) {
	dbVolSlice := make([]autogen_bases.Volumes, len(vs))
	for i, v := range vs {
		dbVol, err := VolumeToDB(v)
		if err != nil {
			return dbVolSlice, err
		}
		dbVolSlice[i] = dbVol
	}
	return dbVolSlice, nil
}

func VolumeToDB(v v1.Volume) (autogen_bases.Volumes, error) {
	dbVolume := autogen_bases.Volumes{
		VolumeID:             0,
		VolumeName:           v.Name,
		VolumeKeyValuesJSONb: "",
	}
	bytesArray, err := json.Marshal(v)
	if err != nil {
		return dbVolume, err
	}
	dbVolume.VolumeKeyValuesJSONb = string(bytesArray)
	return dbVolume, err
}
