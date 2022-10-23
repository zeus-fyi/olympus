package autogen_bases

type ContainersPorts struct {
	ChartSubcomponentChildClassTypeID int `db:"chart_subcomponent_child_class_type_id"`
	ContainerID                       int `db:"container_id"`
	PortID                            int `db:"port_id"`
}
type ContainersPortsSlice []ContainersPorts

func (c *ContainersPorts) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ChartSubcomponentChildClassTypeID, c.ContainerID, c.PortID}
	}
	return pgValues
}
func (c *ContainersPorts) GetTableColumns() (columnValues []string) {
	columnValues = []string{"chart_subcomponent_child_class_type_id", "container_id", "port_id"}
	return columnValues
}
func (c *ContainersPorts) GetTableName() (tableName string) {
	tableName = "containers_ports"
	return tableName
}
