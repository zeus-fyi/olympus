package models

type ContainerComputeResources struct {
	ComputeResourcesID             int    `db:"compute_resources_id"`
	ComputeResourcesKeyValuesJSONb string `db:"compute_resources_key_values_jsonb"`
}
