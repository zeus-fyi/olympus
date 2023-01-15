package hestia_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type KeyTypes struct {
	KeyTypeID   int    `db:"key_type_id" json:"key_type_id"`
	KeyTypeName string `db:"key_type_name" json:"key_type_name"`
}
type KeyTypesSlice []KeyTypes

func (k *KeyTypes) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{k.KeyTypeID, k.KeyTypeName}
	}
	return pgValues
}
func (k *KeyTypes) GetTableColumns() (columnValues []string) {
	columnValues = []string{"key_type_id", "key_type_name"}
	return columnValues
}
func (k *KeyTypes) GetTableName() (tableName string) {
	tableName = "key_types"
	return tableName
}
