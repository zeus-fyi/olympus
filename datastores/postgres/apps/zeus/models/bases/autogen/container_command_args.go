package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ContainerCommandArgs struct {
	ArgsValues    string `db:"args_values" json:"args_values"`
	CommandArgsID int    `db:"command_args_id" json:"command_args_id"`
	CommandValues string `db:"command_values" json:"command_values"`
}
type ContainerCommandArgsSlice []ContainerCommandArgs

func (c *ContainerCommandArgs) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ArgsValues, c.CommandArgsID, c.CommandValues}
	}
	return pgValues
}
func (c *ContainerCommandArgs) GetTableColumns() (columnValues []string) {
	columnValues = []string{"args_values", "command_args_id", "command_values"}
	return columnValues
}
func (c *ContainerCommandArgs) GetTableName() (tableName string) {
	tableName = "container_command_args"
	return tableName
}
