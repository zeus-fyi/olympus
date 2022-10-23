package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ContainerComputeResources struct {
	ComputeResourcesKeyValuesJSONb string `db:"compute_resources_key_values_jsonb"`
	ComputeResourcesID             int    `db:"compute_resources_id"`
}
type ContainerComputeResourcesSlice []ContainerComputeResources

func (c *ContainerComputeResources) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ComputeResourcesKeyValuesJSONb, c.ComputeResourcesID}
	}
	return pgValues
}
func (c *ContainerComputeResources) GetTableColumns() (columnValues []string) {
	columnValues = []string{"compute_resources_key_values_jsonb", "compute_resources_id"}
	return columnValues
}
func (c *ContainerComputeResources) GetTableName() (tableName string) {
	tableName = "container_compute_resources"
	return tableName
}
