package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ContainersSecurityContext struct {
	ContainerSecurityContextID int `db:"container_security_context_id" json:"containerSecurityContextID"`
	ContainerID                int `db:"container_id" json:"containerID"`
}
type ContainersSecurityContextSlice []ContainersSecurityContext

func (c *ContainersSecurityContext) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ContainerSecurityContextID, c.ContainerID}
	}
	return pgValues
}
func (c *ContainersSecurityContext) GetTableColumns() (columnValues []string) {
	columnValues = []string{"container_security_context_id", "container_id"}
	return columnValues
}
func (c *ContainersSecurityContext) GetTableName() (tableName string) {
	tableName = "containers_security_context"
	return tableName
}
