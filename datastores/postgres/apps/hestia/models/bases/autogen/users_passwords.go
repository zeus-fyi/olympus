package hestia_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type UsersPasswords struct {
	UserID     int    `db:"user_id" json:"userID"`
	Password   string `db:"password" json:"password"`
	PasswordID int    `db:"password_id" json:"passwordID"`
}
type UsersLoginsSlice []UsersPasswords

func (u *UsersPasswords) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{u.UserID, u.Password, u.PasswordID}
	}
	return pgValues
}
func (u *UsersPasswords) GetTableColumns() (columnValues []string) {
	columnValues = []string{"user_id", "password", "password_id"}
	return columnValues
}
func (u *UsersPasswords) GetTableName() (tableName string) {
	tableName = "users_passwords"
	return tableName
}
