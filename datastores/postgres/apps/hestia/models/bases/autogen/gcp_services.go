package hestia_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type GcpServices struct {
	ServiceID          string `db:"service_id" json:"serviceID"`
	DisplayName        string `db:"display_name" json:"displayName"`
	BusinessEntityName string `db:"business_entity_name" json:"businessEntityName"`
	Name               string `db:"name" json:"name"`
}
type GcpServicesSlice []GcpServices

func (g *GcpServices) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{g.ServiceID, g.DisplayName, g.BusinessEntityName, g.Name}
	}
	return pgValues
}
func (g *GcpServices) GetTableColumns() (columnValues []string) {
	columnValues = []string{"service_id", "display_name", "business_entity_name", "name"}
	return columnValues
}
func (g *GcpServices) GetTableName() (tableName string) {
	tableName = "gcp_services"
	return tableName
}
