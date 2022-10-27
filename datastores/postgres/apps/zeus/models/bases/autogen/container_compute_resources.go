package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ContainerComputeResources struct {
	ComputeResourcesID             int    `db:"compute_resources_id" json:"compute_resources_id"`
	ComputeResourcesKeyValuesJSONb string `db:"compute_resources_key_values_jsonb" json:"compute_resources_key_values_jsonb"`
}
type ContainerComputeResourcesSlice []ContainerComputeResources

func (c *ContainerComputeResources) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ComputeResourcesID, c.ComputeResourcesKeyValuesJSONb}
	}
	return pgValues
}
func (c *ContainerComputeResources) GetTableColumns() (columnValues []string) {
	columnValues = []string{"compute_resources_id", "compute_resources_key_values_jsonb"}
	return columnValues
}
func (c *ContainerComputeResources) GetTableName() (tableName string) {
	tableName = "container_compute_resources"
	return tableName
}
