package hestia_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type UsersKeyServices struct {
	ServiceID int    `db:"service_id" json:"serviceID"`
	PublicKey string `db:"public_key" json:"publicKey"`
}
type UsersKeyServicesSlice []UsersKeyServices

func (u *UsersKeyServices) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{u.ServiceID, u.PublicKey}
	}
	return pgValues
}
func (u *UsersKeyServices) GetTableColumns() (columnValues []string) {
	columnValues = []string{"service_id", "public_key"}
	return columnValues
}
func (u *UsersKeyServices) GetTableName() (tableName string) {
	tableName = "users_key_services"
	return tableName
}
