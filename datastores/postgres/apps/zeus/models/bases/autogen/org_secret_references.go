package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type OrgSecretReferences struct {
	OrgID      int    `db:"org_id" json:"orgID"`
	SecretID   int    `db:"secret_id" json:"secretID"`
	SecretName string `db:"secret_name" json:"secretName"`
}
type OrgSecretReferencesSlice []OrgSecretReferences

func (o *OrgSecretReferences) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{o.OrgID, o.SecretID, o.SecretName}
	}
	return pgValues
}
func (o *OrgSecretReferences) GetTableColumns() (columnValues []string) {
	columnValues = []string{"org_id", "secret_id", "secret_name"}
	return columnValues
}
func (o *OrgSecretReferences) GetTableName() (tableName string) {
	tableName = "org_secret_references"
	return tableName
}
