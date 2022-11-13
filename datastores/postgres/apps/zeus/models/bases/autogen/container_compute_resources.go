package autogen_bases

type ContainerComputeResources struct {
	ComputeResourcesCpuLimit                string `db:"compute_resources_cpu_limit" json:"computeResourcesCpuLimit"`
	ComputeResourcesRamRequest              string `db:"compute_resources_ram_request" json:"computeResourcesRamRequest"`
	ComputeResourcesRamLimit                string `db:"compute_resources_ram_limit" json:"computeResourcesRamLimit"`
	ComputeResourcesEphemeralStorageRequest string `db:"compute_resources_ephemeral_storage_request" json:"computeResourcesEphemeralStorageRequest"`
	ComputeResourcesEphemeralStorageLimit   string `db:"compute_resources_ephemeral_storage_limit" json:"computeResourcesEphemeralStorageLimit"`
	ComputeResourcesID                      int    `db:"compute_resources_id" json:"computeResourcesID"`
	ComputeResourcesCpuRequest              string `db:"compute_resources_cpu_request" json:"computeResourcesCpuRequest"`
}
type ContainerComputeResourcesSlice []ContainerComputeResources

func (c *ContainerComputeResources) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ComputeResourcesCpuLimit, c.ComputeResourcesRamRequest, c.ComputeResourcesRamLimit, c.ComputeResourcesEphemeralStorageRequest, c.ComputeResourcesEphemeralStorageLimit, c.ComputeResourcesID, c.ComputeResourcesCpuRequest}
	}
	return pgValues
}
func (c *ContainerComputeResources) GetTableColumns() (columnValues []string) {
	columnValues = []string{"compute_resources_cpu_limit", "compute_resources_ram_request", "compute_resources_ram_limit", "compute_resources_ephemeral_storage_request", "compute_resources_ephemeral_storage_limit", "compute_resources_id", "compute_resources_cpu_request"}
	return columnValues
}
func (c *ContainerComputeResources) GetTableName() (tableName string) {
	tableName = "container_compute_resources"
	return tableName
}
