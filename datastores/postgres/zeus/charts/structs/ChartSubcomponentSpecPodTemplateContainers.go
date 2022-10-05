package models

type ChartSubcomponentSpecPodTemplateContainers struct {
	ChartSubcomponentChildClassTypeID int  `db:"chart_subcomponent_child_class_type_id"`
	ContainerID                       int  `db:"container_id"`
	IsInitContainer                   bool `db:"is_init_container"`
	ContainerSortOrder                int  `db:"container_sort_order"`
}
