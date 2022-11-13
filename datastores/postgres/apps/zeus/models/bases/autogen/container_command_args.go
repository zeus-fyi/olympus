package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ContainerCommandArgs struct {
	CommandArgsID int    `db:"command_args_id" json:"commandArgsID"`
	CommandValues string `db:"command_values" json:"commandValues"`
	ArgsValues    string `db:"args_values" json:"argsValues"`
}
type ContainerCommandArgsSlice []ContainerCommandArgs

func (c *ContainerCommandArgs) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.CommandArgsID, c.CommandValues, c.ArgsValues}
	}
	return pgValues
}
func (c *ContainerCommandArgs) GetTableColumns() (columnValues []string) {
	columnValues = []string{"command_args_id", "command_values", "args_values"}
	return columnValues
}
func (c *ContainerCommandArgs) GetTableName() (tableName string) {
	tableName = "container_command_args"
	return tableName
}
