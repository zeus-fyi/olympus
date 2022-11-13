package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ContainersComputeResources struct {
	ComputeResourcesID int `db:"compute_resources_id" json:"computeResourcesID"`
	ContainerID        int `db:"container_id" json:"containerID"`
}
type ContainersComputeResourcesSlice []ContainersComputeResources

func (c *ContainersComputeResources) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ComputeResourcesID, c.ContainerID}
	}
	return pgValues
}
func (c *ContainersComputeResources) GetTableColumns() (columnValues []string) {
	columnValues = []string{"compute_resources_id", "container_id"}
	return columnValues
}
func (c *ContainersComputeResources) GetTableName() (tableName string) {
	tableName = "containers_compute_resources"
	return tableName
}
