package hestia_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type Resources struct {
	ResourceID int    `db:"resource_id" json:"resourceID"`
	Type       string `db:"type" json:"type"`
}
type ResourcesSlice []Resources

func (r *Resources) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{r.ResourceID, r.Type}
	}
	return pgValues
}
func (r *Resources) GetTableColumns() (columnValues []string) {
	columnValues = []string{"resource_id", "type"}
	return columnValues
}
func (r *Resources) GetTableName() (tableName string) {
	tableName = "resources"
	return tableName
}
