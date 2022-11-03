package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ContainersComputeResources struct {
	ContainerID        int `db:"container_id" json:"container_id"`
	ComputeResourcesID int `db:"compute_resources_id" json:"compute_resources_id"`
}
type ContainersComputeResourcesSlice []ContainersComputeResources

func (c *ContainersComputeResources) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ContainerID, c.ComputeResourcesID}
	}
	return pgValues
}
func (c *ContainersComputeResources) GetTableColumns() (columnValues []string) {
	columnValues = []string{"container_id", "compute_resources_id"}
	return columnValues
}
func (c *ContainersComputeResources) GetTableName() (tableName string) {
	tableName = "containers_compute_resources"
	return tableName
}
