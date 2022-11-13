package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ContainerVolumeMounts struct {
	VolumeMountID   int    `db:"volume_mount_id" json:"volumeMountID"`
	VolumeMountPath string `db:"volume_mount_path" json:"volumeMountPath"`
	VolumeName      string `db:"volume_name" json:"volumeName"`
	VolumeReadOnly  bool   `db:"volume_read_only" json:"volumeReadOnly"`
	VolumeSubPath   string `db:"volume_sub_path" json:"volumeSubPath"`
}
type ContainerVolumeMountsSlice []ContainerVolumeMounts

func (c *ContainerVolumeMounts) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.VolumeMountID, c.VolumeMountPath, c.VolumeName, c.VolumeReadOnly, c.VolumeSubPath}
	}
	return pgValues
}
func (c *ContainerVolumeMounts) GetTableColumns() (columnValues []string) {
	columnValues = []string{"volume_mount_id", "volume_mount_path", "volume_name", "volume_read_only", "volume_sub_path"}
	return columnValues
}
func (c *ContainerVolumeMounts) GetTableName() (tableName string) {
	tableName = "container_volume_mounts"
	return tableName
}
