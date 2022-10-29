package autogen_bases

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

type UsersKeyGroups struct {
	KeyGroupID   int       `db:"key_group_id" json:"key_group_id"`
	KeyGroupName string    `db:"key_group_name" json:"key_group_name"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
	UserID       int       `db:"user_id" json:"user_id"`
	PublicKey    string    `db:"public_key" json:"public_key"`
}
type UsersKeyGroupsSlice []UsersKeyGroups

func (u *UsersKeyGroups) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{u.KeyGroupID, u.KeyGroupName, u.UpdatedAt, u.UserID, u.PublicKey}
	}
	return pgValues
}
func (u *UsersKeyGroups) GetTableColumns() (columnValues []string) {
	columnValues = []string{"key_group_id", "key_group_name", "updated_at", "user_id", "public_key"}
	return columnValues
}
func (u *UsersKeyGroups) GetTableName() (tableName string) {
	tableName = "users_key_groups"
	return tableName
}
