package hestia_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type Services struct {
	ServiceID   int    `db:"service_id" json:"serviceID"`
	ServiceName string `db:"service_name" json:"serviceName"`
}
type ServicesSlice []Services

func (s *Services) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{s.ServiceID, s.ServiceName}
	}
	return pgValues
}
func (s *Services) GetTableColumns() (columnValues []string) {
	columnValues = []string{"service_id", "service_name"}
	return columnValues
}
func (s *Services) GetTableName() (tableName string) {
	tableName = "services"
	return tableName
}
