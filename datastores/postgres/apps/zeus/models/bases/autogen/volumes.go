package autogen_bases

type Volumes struct {
	VolumeName           string `db:"volume_name" json:"volumeName"`
	VolumeKeyValuesJSONb string `db:"volume_key_values_jsonb" json:"volumeKeyValuesJsonb"`
	VolumeID             int    `db:"volume_id" json:"volumeID"`
}
type VolumesSlice []Volumes

func (v *Volumes) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{v.VolumeName, v.VolumeKeyValuesJSONb, v.VolumeID}
	}
	return pgValues
}
func (v *Volumes) GetTableColumns() (columnValues []string) {
	columnValues = []string{"volume_name", "volume_key_values_jsonb", "volume_id"}
	return columnValues
}
func (v *Volumes) GetTableName() (tableName string) {
	tableName = "volumes"
	return tableName
}
