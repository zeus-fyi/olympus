package autogen_bases

type ContainersVolumeMounts struct {
	VolumeMountID                     int `db:"volume_mount_id"`
	ChartSubcomponentChildClassTypeID int `db:"chart_subcomponent_child_class_type_id"`
	ContainerID                       int `db:"container_id"`
}
type ContainersVolumeMountsSlice []ContainersVolumeMounts

func (c *ContainersVolumeMounts) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.VolumeMountID, c.ChartSubcomponentChildClassTypeID, c.ContainerID}
	}
	return pgValues
}
func (c *ContainersVolumeMounts) GetTableColumns() (columnValues []string) {
	columnValues = []string{"volume_mount_id", "chart_subcomponent_child_class_type_id", "container_id"}
	return columnValues
}
func (c *ContainersVolumeMounts) GetTableName() (tableName string) {
	tableName = "containers_volume_mounts"
	return tableName
}
