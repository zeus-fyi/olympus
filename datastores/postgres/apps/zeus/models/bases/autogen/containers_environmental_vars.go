package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ContainersEnvironmentalVars struct {
	EnvID                             int `db:"env_id" json:"env_id"`
	ChartSubcomponentChildClassTypeID int `db:"chart_subcomponent_child_class_type_id" json:"chart_subcomponent_child_class_type_id"`
	ContainerID                       int `db:"container_id" json:"container_id"`
}
type ContainersEnvironmentalVarsSlice []ContainersEnvironmentalVars

func (c *ContainersEnvironmentalVars) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.EnvID, c.ChartSubcomponentChildClassTypeID, c.ContainerID}
	}
	return pgValues
}
func (c *ContainersEnvironmentalVars) GetTableColumns() (columnValues []string) {
	columnValues = []string{"env_id", "chart_subcomponent_child_class_type_id", "container_id"}
	return columnValues
}
func (c *ContainersEnvironmentalVars) GetTableName() (tableName string) {
	tableName = "containers_environmental_vars"
	return tableName
}
