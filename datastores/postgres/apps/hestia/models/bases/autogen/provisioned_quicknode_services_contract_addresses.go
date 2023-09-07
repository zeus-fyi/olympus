package hestia_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ProvisionedQuickNodeServicesContractAddresses struct {
	QuickNodeID     string `db:"quicknode_id" json:"quicknodeID"`
	ContractAddress string `db:"contract_address" json:"contractAddress"`
}
type ProvisionedQuickNodeServicesContractAddressesSlice []ProvisionedQuickNodeServicesContractAddresses

func (p *ProvisionedQuickNodeServicesContractAddresses) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{p.QuickNodeID, p.ContractAddress}
	}
	return pgValues
}
func (p *ProvisionedQuickNodeServicesContractAddresses) GetTableColumns() (columnValues []string) {
	columnValues = []string{"quicknode_id", "contract_address"}
	return columnValues
}
func (p *ProvisionedQuickNodeServicesContractAddresses) GetTableName() (tableName string) {
	tableName = "provisioned_quicknode_services_contract_addresses"
	return tableName
}
