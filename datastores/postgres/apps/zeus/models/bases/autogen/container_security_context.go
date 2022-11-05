package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ContainerSecurityContext struct {
	ContainerSecurityContextID int    `db:"container_security_context_id" json:"container_security_context_id"`
	SecurityContextKeyValues   string `db:"security_context_key_values" json:"security_context_key_values"`
}
type ContainerSecurityContextSlice []ContainerSecurityContext

func (c *ContainerSecurityContext) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ContainerSecurityContextID, c.SecurityContextKeyValues}
	}
	return pgValues
}
func (c *ContainerSecurityContext) GetTableColumns() (columnValues []string) {
	columnValues = []string{"container_security_context_id", "security_context_key_values"}
	return columnValues
}
func (c *ContainerSecurityContext) GetTableName() (tableName string) {
	tableName = "container_security_context"
	return tableName
}
