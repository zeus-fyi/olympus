package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ContainersVolumes struct {
	ChartSubcomponentChildClassTypeID int `db:"chart_subcomponent_child_class_type_id" json:"chartSubcomponentChildClassTypeID"`
	VolumeID                          int `db:"volume_id" json:"volumeID"`
}
type ContainersVolumesSlice []ContainersVolumes

func (c *ContainersVolumes) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ChartSubcomponentChildClassTypeID, c.VolumeID}
	}
	return pgValues
}
func (c *ContainersVolumes) GetTableColumns() (columnValues []string) {
	columnValues = []string{"chart_subcomponent_child_class_type_id", "volume_id"}
	return columnValues
}
func (c *ContainersVolumes) GetTableName() (tableName string) {
	tableName = "containers_volumes"
	return tableName
}
