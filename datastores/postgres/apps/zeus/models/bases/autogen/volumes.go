package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type Volumes struct {
	VolumeID             int    `db:"volume_id" json:"volume_id"`
	VolumeName           string `db:"volume_name" json:"volume_name"`
	VolumeKeyValuesJSONb string `db:"volume_key_values_jsonb" json:"volume_key_values_jsonb"`
}
type VolumesSlice []Volumes

func (v *Volumes) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{v.VolumeID, v.VolumeName, v.VolumeKeyValuesJSONb}
	}
	return pgValues
}
func (v *Volumes) GetTableColumns() (columnValues []string) {
	columnValues = []string{"volume_id", "volume_name", "volume_key_values_jsonb"}
	return columnValues
}
func (v *Volumes) GetTableName() (tableName string) {
	tableName = "volumes"
	return tableName
}
