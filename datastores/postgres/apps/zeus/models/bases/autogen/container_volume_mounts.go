package autogen_bases

type ContainerVolumeMounts struct {
	VolumeMountID   int    `db:"volume_mount_id"`
	VolumeMountPath string `db:"volume_mount_path"`
	VolumeName      string `db:"volume_name"`
}
type ContainerVolumeMountsSlice []ContainerVolumeMounts

func (c *ContainerVolumeMounts) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.VolumeMountID, c.VolumeMountPath, c.VolumeName}
	}
	return pgValues
}
func (c *ContainerVolumeMounts) GetTableColumns() (columnValues []string) {
	columnValues = []string{"volume_mount_id", "volume_mount_path", "volume_name"}
	return columnValues
}
func (c *ContainerVolumeMounts) GetTableName() (tableName string) {
	tableName = "container_volume_mounts"
	return tableName
}
