package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ChartSubcomponentSpecPodTemplateContainers struct {
	IsInitContainer                   bool `db:"is_init_container"`
	ContainerSortOrder                int  `db:"container_sort_order"`
	ChartSubcomponentChildClassTypeID int  `db:"chart_subcomponent_child_class_type_id"`
	ContainerID                       int  `db:"container_id"`
}
type ChartSubcomponentSpecPodTemplateContainersSlice []ChartSubcomponentSpecPodTemplateContainers

func (c *ChartSubcomponentSpecPodTemplateContainers) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.IsInitContainer, c.ContainerSortOrder, c.ChartSubcomponentChildClassTypeID, c.ContainerID}
	}
	return pgValues
}
func (c *ChartSubcomponentSpecPodTemplateContainers) GetTableColumns() (columnValues []string) {
	columnValues = []string{"is_init_container", "container_sort_order", "chart_subcomponent_child_class_type_id", "container_id"}
	return columnValues
}
func (c *ChartSubcomponentSpecPodTemplateContainers) GetTableName() (tableName string) {
	tableName = "chart_subcomponent_spec_pod_template_containers"
	return tableName
}
