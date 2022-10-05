package models

type ContainersComputeResources struct {
	ComputeResourcesID int `db:"compute_resources_id"`
	ContainerID        int `db:"container_id"`
}
