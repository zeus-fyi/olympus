package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ContainerComputeResources struct {
	ComputeResourcesRamRequest              string `db:"compute_resources_ram_request" json:"compute_resources_ram_request"`
	ComputeResourcesRamLimit                string `db:"compute_resources_ram_limit" json:"compute_resources_ram_limit"`
	ComputeResourcesEphemeralStorageRequest string `db:"compute_resources_ephemeral_storage_request" json:"compute_resources_ephemeral_storage_request"`
	ComputeResourcesEphemeralStorageLimit   string `db:"compute_resources_ephemeral_storage_limit" json:"compute_resources_ephemeral_storage_limit"`
	ComputeResourcesID                      int    `db:"compute_resources_id" json:"compute_resources_id"`
	ComputeResourcesCpuRequest              string `db:"compute_resources_cpu_request" json:"compute_resources_cpu_request"`
	ComputeResourcesCpuLimit                string `db:"compute_resources_cpu_limit" json:"compute_resources_cpu_limit"`
}
type ContainerComputeResourcesSlice []ContainerComputeResources

func (c *ContainerComputeResources) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ComputeResourcesRamRequest, c.ComputeResourcesRamLimit, c.ComputeResourcesEphemeralStorageRequest, c.ComputeResourcesEphemeralStorageLimit, c.ComputeResourcesID, c.ComputeResourcesCpuRequest, c.ComputeResourcesCpuLimit}
	}
	return pgValues
}
func (c *ContainerComputeResources) GetTableColumns() (columnValues []string) {
	columnValues = []string{"compute_resources_ram_request", "compute_resources_ram_limit", "compute_resources_ephemeral_storage_request", "compute_resources_ephemeral_storage_limit", "compute_resources_id", "compute_resources_cpu_request", "compute_resources_cpu_limit"}
	return columnValues
}
func (c *ContainerComputeResources) GetTableName() (tableName string) {
	tableName = "container_compute_resources"
	return tableName
}
