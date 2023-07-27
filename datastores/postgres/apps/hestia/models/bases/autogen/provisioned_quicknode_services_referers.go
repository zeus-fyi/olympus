package hestia_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ProvisionedQuicknodeServicesReferers struct {
	QuicknodeID string `db:"quicknode_id" json:"quicknodeID"`
	Referer     string `db:"referer" json:"referer"`
}
type ProvisionedQuicknodeServicesReferersSlice []ProvisionedQuicknodeServicesReferers

func (p *ProvisionedQuicknodeServicesReferers) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{p.QuicknodeID, p.Referer}
	}
	return pgValues
}
func (p *ProvisionedQuicknodeServicesReferers) GetTableColumns() (columnValues []string) {
	columnValues = []string{"quicknode_id", "referer"}
	return columnValues
}
func (p *ProvisionedQuicknodeServicesReferers) GetTableName() (tableName string) {
	tableName = "provisioned_quicknode_services_referers"
	return tableName
}
