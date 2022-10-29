package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type UserOrgs struct {
	OrgID  int `db:"org_id" json:"org_id"`
	UserID int `db:"user_id" json:"user_id"`
}
type UserOrgsSlice []UserOrgs

func (u *UserOrgs) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{u.OrgID, u.UserID}
	}
	return pgValues
}
func (u *UserOrgs) GetTableColumns() (columnValues []string) {
	columnValues = []string{"org_id", "user_id"}
	return columnValues
}
func (u *UserOrgs) GetTableName() (tableName string) {
	tableName = "user_orgs"
	return tableName
}
