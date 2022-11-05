package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ContainerVolumeMounts struct {
	VolumeReadOnly  bool   `db:"volume_read_only" json:"volume_read_only"`
	VolumeSubPath   string `db:"volume_sub_path" json:"volume_sub_path"`
	VolumeMountID   int    `db:"volume_mount_id" json:"volume_mount_id"`
	VolumeMountPath string `db:"volume_mount_path" json:"volume_mount_path"`
	VolumeName      string `db:"volume_name" json:"volume_name"`
}
type ContainerVolumeMountsSlice []ContainerVolumeMounts

func (c *ContainerVolumeMounts) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.VolumeReadOnly, c.VolumeSubPath, c.VolumeMountID, c.VolumeMountPath, c.VolumeName}
	}
	return pgValues
}
func (c *ContainerVolumeMounts) GetTableColumns() (columnValues []string) {
	columnValues = []string{"volume_read_only", "volume_sub_path", "volume_mount_id", "volume_mount_path", "volume_name"}
	return columnValues
}
func (c *ContainerVolumeMounts) GetTableName() (tableName string) {
	tableName = "container_volume_mounts"
	return tableName
}
