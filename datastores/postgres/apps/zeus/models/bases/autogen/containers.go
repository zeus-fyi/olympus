package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type Containers struct {
	ContainerImageID         string `db:"container_image_id" json:"containerImageID"`
	ContainerVersionTag      string `db:"container_version_tag" json:"containerVersionTag"`
	ContainerPlatformOs      string `db:"container_platform_os" json:"containerPlatformOs"`
	ContainerRepository      string `db:"container_repository" json:"containerRepository"`
	ContainerImagePullPolicy string `db:"container_image_pull_policy" json:"containerImagePullPolicy"`
	IsInitContainer          bool   `db:"is_init_container" json:"isInitContainer"`
	ContainerID              int    `db:"container_id" json:"containerID"`
	ContainerName            string `db:"container_name" json:"containerName"`
}
type ContainersSlice []Containers

func (c *Containers) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ContainerImageID, c.ContainerVersionTag, c.ContainerPlatformOs, c.ContainerRepository, c.ContainerImagePullPolicy, c.IsInitContainer, c.ContainerID, c.ContainerName}
	}
	return pgValues
}
func (c *Containers) GetTableColumns() (columnValues []string) {
	columnValues = []string{"container_image_id", "container_version_tag", "container_platform_os", "container_repository", "container_image_pull_policy", "is_init_container", "container_id", "container_name"}
	return columnValues
}
func (c *Containers) GetTableName() (tableName string) {
	tableName = "containers"
	return tableName
}
