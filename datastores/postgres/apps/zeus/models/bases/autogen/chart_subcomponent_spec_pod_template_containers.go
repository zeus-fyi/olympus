package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ChartSubcomponentSpecPodTemplateContainers struct {
	ChartSubcomponentChildClassTypeID int  `db:"chart_subcomponent_child_class_type_id"`
	ContainerID                       int  `db:"container_id"`
	IsInitContainer                   bool `db:"is_init_container"`
	ContainerSortOrder                int  `db:"container_sort_order"`
}
type ChartSubcomponentSpecPodTemplateContainersSlice []ChartSubcomponentSpecPodTemplateContainers

func (c *ChartSubcomponentSpecPodTemplateContainers) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ChartSubcomponentChildClassTypeID, c.ContainerID, c.IsInitContainer, c.ContainerSortOrder}
	}
	return pgValues
}
func (c *ChartSubcomponentSpecPodTemplateContainers) GetTableColumns() (columnValues []string) {
	columnValues = []string{"chart_subcomponent_child_class_type_id", "container_id", "is_init_container", "container_sort_order"}
	return columnValues
}
func (c *ChartSubcomponentSpecPodTemplateContainers) GetTableName() (tableName string) {
	tableName = "chart_subcomponent_spec_pod_template_containers"
	return tableName
}
