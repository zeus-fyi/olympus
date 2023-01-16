package hestia_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ValidatorsServiceOrgGroupsCloudCtxNs struct {
	CloudCtxNsID int    `db:"cloud_ctx_ns_id" json:"cloudCtxNsID"`
	Pubkey       string `db:"pubkey" json:"pubkey"`
}
type ValidatorsServiceOrgGroupsCloudCtxNsSlice []ValidatorsServiceOrgGroupsCloudCtxNs

func (v *ValidatorsServiceOrgGroupsCloudCtxNs) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{v.CloudCtxNsID, v.Pubkey}
	}
	return pgValues
}
func (v *ValidatorsServiceOrgGroupsCloudCtxNs) GetTableColumns() (columnValues []string) {
	columnValues = []string{"cloud_ctx_ns_id", "pubkey"}
	return columnValues
}
func (v *ValidatorsServiceOrgGroupsCloudCtxNs) GetTableName() (tableName string) {
	tableName = "validators_service_org_groups_cloud_ctx_ns"
	return tableName
}
