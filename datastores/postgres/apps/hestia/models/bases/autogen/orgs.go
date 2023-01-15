package hestia_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type Orgs struct {
	OrgID    int    `db:"org_id" json:"orgID"`
	Metadata string `db:"metadata" json:"metadata"`
}
type OrgsSlice []Orgs

func (o *Orgs) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{o.OrgID, o.Metadata}
	}
	return pgValues
}
func (o *Orgs) GetTableColumns() (columnValues []string) {
	columnValues = []string{"org_id", "metadata"}
	return columnValues
}
func (o *Orgs) GetTableName() (tableName string) {
	tableName = "orgs"
	return tableName
}
