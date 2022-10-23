package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type Containers struct {
	ContainerImagePullPolicy string `db:"container_image_pull_policy"`
	ContainerID              int    `db:"container_id"`
	ContainerName            string `db:"container_name"`
	ContainerImageID         string `db:"container_image_id"`
	ContainerVersionTag      string `db:"container_version_tag"`
	ContainerPlatformOs      string `db:"container_platform_os"`
	ContainerRepository      string `db:"container_repository"`
}
type ContainersSlice []Containers

func (c *Containers) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ContainerImagePullPolicy, c.ContainerID, c.ContainerName, c.ContainerImageID, c.ContainerVersionTag, c.ContainerPlatformOs, c.ContainerRepository}
	}
	return pgValues
}
func (c *Containers) GetTableColumns() (columnValues []string) {
	columnValues = []string{"container_image_pull_policy", "container_id", "container_name", "container_image_id", "container_version_tag", "container_platform_os", "container_repository"}
	return columnValues
}
func (c *Containers) GetTableName() (tableName string) {
	tableName = "containers"
	return tableName
}
