package models

type Volumes struct {
	VolumeID             int    `db:"volume_id"`
	VolumeName           string `db:"volume_name"`
	VolumeKeyValuesJSONb string `db:"volume_key_values_jsonb"`
}
