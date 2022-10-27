package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ChartSubcomponentSpecPodTemplateContainers struct {
	ContainerID                       int  `db:"container_id" json:"container_id"`
	IsInitContainer                   bool `db:"is_init_container" json:"is_init_container"`
	ContainerSortOrder                int  `db:"container_sort_order" json:"container_sort_order"`
	ChartSubcomponentChildClassTypeID int  `db:"chart_subcomponent_child_class_type_id" json:"chart_subcomponent_child_class_type_id"`
}
type ChartSubcomponentSpecPodTemplateContainersSlice []ChartSubcomponentSpecPodTemplateContainers

func (c *ChartSubcomponentSpecPodTemplateContainers) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ContainerID, c.IsInitContainer, c.ContainerSortOrder, c.ChartSubcomponentChildClassTypeID}
	}
	return pgValues
}
func (c *ChartSubcomponentSpecPodTemplateContainers) GetTableColumns() (columnValues []string) {
	columnValues = []string{"container_id", "is_init_container", "container_sort_order", "chart_subcomponent_child_class_type_id"}
	return columnValues
}
func (c *ChartSubcomponentSpecPodTemplateContainers) GetTableName() (tableName string) {
	tableName = "chart_subcomponent_spec_pod_template_containers"
	return tableName
}
