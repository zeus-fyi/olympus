package models

type ContainersVolumeMounts struct {
	ChartSubcomponentChildClassTypeID int `db:"chart_subcomponent_child_class_type_id"`
	ContainerID                       int `db:"container_id"`
	VolumeMountID                     int `db:"volume_mount_id"`
}
