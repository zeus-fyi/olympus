package hestia_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type Users struct {
	UserID   int    `db:"user_id" json:"userID"`
	Metadata string `db:"metadata" json:"metadata"`
}
type UsersSlice []Users

func (u *Users) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{u.UserID, u.Metadata}
	}
	return pgValues
}
func (u *Users) GetTableColumns() (columnValues []string) {
	columnValues = []string{"user_id", "metadata"}
	return columnValues
}
func (u *Users) GetTableName() (tableName string) {
	tableName = "users"
	return tableName
}
