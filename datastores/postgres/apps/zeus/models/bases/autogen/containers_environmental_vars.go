package autogen_bases

type ContainersEnvironmentalVars struct {
	ContainerID                       int `db:"container_id"`
	EnvID                             int `db:"env_id"`
	ChartSubcomponentChildClassTypeID int `db:"chart_subcomponent_child_class_type_id"`
}
type ContainersEnvironmentalVarsSlice []ContainersEnvironmentalVars

func (c *ContainersEnvironmentalVars) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ContainerID, c.EnvID, c.ChartSubcomponentChildClassTypeID}
	}
	return pgValues
}
func (c *ContainersEnvironmentalVars) GetTableColumns() (columnValues []string) {
	columnValues = []string{"container_id", "env_id", "chart_subcomponent_child_class_type_id"}
	return columnValues
}
func (c *ContainersEnvironmentalVars) GetTableName() (tableName string) {
	tableName = "containers_environmental_vars"
	return tableName
}
