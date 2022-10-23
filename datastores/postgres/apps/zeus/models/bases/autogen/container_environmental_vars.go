package autogen_bases

type ContainerEnvironmentalVars struct {
	EnvID int    `db:"env_id"`
	Name  string `db:"name"`
	Value string `db:"value"`
}
type ContainerEnvironmentalVarsSlice []ContainerEnvironmentalVars

func (c *ContainerEnvironmentalVars) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.EnvID, c.Name, c.Value}
	}
	return pgValues
}
func (c *ContainerEnvironmentalVars) GetTableColumns() (columnValues []string) {
	columnValues = []string{"env_id", "name", "value"}
	return columnValues
}
func (c *ContainerEnvironmentalVars) GetTableName() (tableName string) {
	tableName = "container_environmental_vars"
	return tableName
}
