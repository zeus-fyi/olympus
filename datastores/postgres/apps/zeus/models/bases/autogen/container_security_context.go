package autogen_bases

type ContainerSecurityContext struct {
	SecurityContextKeyValues   string `db:"security_context_key_values" json:"securityContextKeyValues"`
	ContainerSecurityContextID int    `db:"container_security_context_id" json:"containerSecurityContextID"`
}
type ContainerSecurityContextSlice []ContainerSecurityContext

func (c *ContainerSecurityContext) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.SecurityContextKeyValues, c.ContainerSecurityContextID}
	}
	return pgValues
}
func (c *ContainerSecurityContext) GetTableColumns() (columnValues []string) {
	columnValues = []string{"security_context_key_values", "container_security_context_id"}
	return columnValues
}
func (c *ContainerSecurityContext) GetTableName() (tableName string) {
	tableName = "container_security_context"
	return tableName
}
