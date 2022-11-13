package autogen_bases

type ContainersVolumeMounts struct {
	ChartSubcomponentChildClassTypeID int `db:"chart_subcomponent_child_class_type_id" json:"chartSubcomponentChildClassTypeID"`
	ContainerID                       int `db:"container_id" json:"containerID"`
	VolumeMountID                     int `db:"volume_mount_id" json:"volumeMountID"`
}
type ContainersVolumeMountsSlice []ContainersVolumeMounts

func (c *ContainersVolumeMounts) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ChartSubcomponentChildClassTypeID, c.ContainerID, c.VolumeMountID}
	}
	return pgValues
}
func (c *ContainersVolumeMounts) GetTableColumns() (columnValues []string) {
	columnValues = []string{"chart_subcomponent_child_class_type_id", "container_id", "volume_mount_id"}
	return columnValues
}
func (c *ContainersVolumeMounts) GetTableName() (tableName string) {
	tableName = "containers_volume_mounts"
	return tableName
}
