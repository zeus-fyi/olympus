package autogen_bases

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

type UsersKeys struct {
	UserID            int       `db:"user_id" json:"user_id"`
	PublicKeyName     string    `db:"public_key_name" json:"public_key_name"`
	PublicKeyVerified bool      `db:"public_key_verified" json:"public_key_verified"`
	PublicKeyTypeID   int       `db:"public_key_type_id" json:"public_key_type_id"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	PublicKey         string    `db:"public_key" json:"public_key"`
}
type UsersKeysSlice []UsersKeys

func (u *UsersKeys) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{u.UserID, u.PublicKeyName, u.PublicKeyVerified, u.PublicKeyTypeID, u.CreatedAt, u.PublicKey}
	}
	return pgValues
}
func (u *UsersKeys) GetTableColumns() (columnValues []string) {
	columnValues = []string{"user_id", "public_key_name", "public_key_verified", "public_key_type_id", "created_at", "public_key"}
	return columnValues
}
func (u *UsersKeys) GetTableName() (tableName string) {
	tableName = "users_keys"
	return tableName
}
