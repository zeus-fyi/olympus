package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type OrgSecretKeyValReferences struct {
	SecretID        int    `db:"secret_id" json:"secretID"`
	SecretEnvVarRef string `db:"secret_env_var_ref" json:"secretEnvVarRef"`
	SecretKeyRef    string `db:"secret_key_ref" json:"secretKeyRef"`
	SecretNameRef   string `db:"secret_name_ref" json:"secretNameRef"`
}
type OrgSecretKeyValReferencesSlice []OrgSecretKeyValReferences

func (o *OrgSecretKeyValReferences) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{o.SecretID, o.SecretEnvVarRef, o.SecretKeyRef, o.SecretNameRef}
	}
	return pgValues
}
func (o *OrgSecretKeyValReferences) GetTableColumns() (columnValues []string) {
	columnValues = []string{"secret_id", "secret_env_var_ref", "secret_key_ref", "secret_name_ref"}
	return columnValues
}
func (o *OrgSecretKeyValReferences) GetTableName() (tableName string) {
	tableName = "org_secret_key_val_references"
	return tableName
}
