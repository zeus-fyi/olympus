package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type Containers struct {
	ContainerName            string `db:"container_name" json:"container_name"`
	ContainerImageID         string `db:"container_image_id" json:"container_image_id"`
	ContainerVersionTag      string `db:"container_version_tag" json:"container_version_tag"`
	ContainerPlatformOs      string `db:"container_platform_os" json:"container_platform_os"`
	ContainerRepository      string `db:"container_repository" json:"container_repository"`
	ContainerImagePullPolicy string `db:"container_image_pull_policy" json:"container_image_pull_policy"`
	ContainerID              int    `db:"container_id" json:"container_id"`
}
type ContainersSlice []Containers

func (c *Containers) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ContainerName, c.ContainerImageID, c.ContainerVersionTag, c.ContainerPlatformOs, c.ContainerRepository, c.ContainerImagePullPolicy, c.ContainerID}
	}
	return pgValues
}
func (c *Containers) GetTableColumns() (columnValues []string) {
	columnValues = []string{"container_name", "container_image_id", "container_version_tag", "container_platform_os", "container_repository", "container_image_pull_policy", "container_id"}
	return columnValues
}
func (c *Containers) GetTableName() (tableName string) {
	tableName = "containers"
	return tableName
}
