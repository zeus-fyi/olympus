package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ContainersPorts struct {
	ContainerID                       int `db:"container_id"`
	PortID                            int `db:"port_id"`
	ChartSubcomponentChildClassTypeID int `db:"chart_subcomponent_child_class_type_id"`
}
type ContainersPortsSlice []ContainersPorts

func (c *ContainersPorts) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ContainerID, c.PortID, c.ChartSubcomponentChildClassTypeID}
	}
	return pgValues
}
func (c *ContainersPorts) GetTableColumns() (columnValues []string) {
	columnValues = []string{"container_id", "port_id", "chart_subcomponent_child_class_type_id"}
	return columnValues
}
func (c *ContainersPorts) GetTableName() (tableName string) {
	tableName = "containers_ports"
	return tableName
}
