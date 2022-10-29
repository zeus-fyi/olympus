package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type OrgUsers struct {
	OrgID  int `db:"org_id" json:"org_id"`
	UserID int `db:"user_id" json:"user_id"`
}
type OrgUsersSlice []OrgUsers

func (o *OrgUsers) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{o.OrgID, o.UserID}
	}
	return pgValues
}
func (o *OrgUsers) GetTableColumns() (columnValues []string) {
	columnValues = []string{"org_id", "user_id"}
	return columnValues
}
func (o *OrgUsers) GetTableName() (tableName string) {
	tableName = "org_users"
	return tableName
}
