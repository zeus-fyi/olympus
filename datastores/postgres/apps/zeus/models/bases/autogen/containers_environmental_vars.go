package autogen_bases

type ContainersEnvironmentalVars struct {
	ChartSubcomponentChildClassTypeID int `db:"chart_subcomponent_child_class_type_id" json:"chartSubcomponentChildClassTypeID"`
	ContainerID                       int `db:"container_id" json:"containerID"`
	EnvID                             int `db:"env_id" json:"envID"`
}
type ContainersEnvironmentalVarsSlice []ContainersEnvironmentalVars

func (c *ContainersEnvironmentalVars) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ChartSubcomponentChildClassTypeID, c.ContainerID, c.EnvID}
	}
	return pgValues
}
func (c *ContainersEnvironmentalVars) GetTableColumns() (columnValues []string) {
	columnValues = []string{"chart_subcomponent_child_class_type_id", "container_id", "env_id"}
	return columnValues
}
func (c *ContainersEnvironmentalVars) GetTableName() (tableName string) {
	tableName = "containers_environmental_vars"
	return tableName
}
