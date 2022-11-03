package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ContainersCommandArgs struct {
	CommandArgsID int `db:"command_args_id" json:"command_args_id"`
	ContainerID   int `db:"container_id" json:"container_id"`
}
type ContainersCommandArgsSlice []ContainersCommandArgs

func (c *ContainersCommandArgs) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.CommandArgsID, c.ContainerID}
	}
	return pgValues
}
func (c *ContainersCommandArgs) GetTableColumns() (columnValues []string) {
	columnValues = []string{"command_args_id", "container_id"}
	return columnValues
}
func (c *ContainersCommandArgs) GetTableName() (tableName string) {
	tableName = "containers_command_args"
	return tableName
}
