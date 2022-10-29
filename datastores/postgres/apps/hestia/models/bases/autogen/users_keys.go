package autogen_bases

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

type UsersKeys struct {
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	PublicKey         string    `db:"public_key" json:"public_key"`
	UserID            int       `db:"user_id" json:"user_id"`
	PublicKeyName     string    `db:"public_key_name" json:"public_key_name"`
	PublicKeyVerified bool      `db:"public_key_verified" json:"public_key_verified"`
	PublicKeyTypeID   int       `db:"public_key_type_id" json:"public_key_type_id"`
}
type UsersKeysSlice []UsersKeys

func (u *UsersKeys) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{u.CreatedAt, u.PublicKey, u.UserID, u.PublicKeyName, u.PublicKeyVerified, u.PublicKeyTypeID}
	}
	return pgValues
}
func (u *UsersKeys) GetTableColumns() (columnValues []string) {
	columnValues = []string{"created_at", "public_key", "user_id", "public_key_name", "public_key_verified", "public_key_type_id"}
	return columnValues
}
func (u *UsersKeys) GetTableName() (tableName string) {
	tableName = "users_keys"
	return tableName
}
