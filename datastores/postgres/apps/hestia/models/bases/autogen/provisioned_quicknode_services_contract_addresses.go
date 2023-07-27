package hestia_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ProvisionedQuicknodeServicesContractAddresses struct {
	QuicknodeID     string `db:"quicknode_id" json:"quicknodeID"`
	ContractAddress string `db:"contract_address" json:"contractAddress"`
}
type ProvisionedQuicknodeServicesContractAddressesSlice []ProvisionedQuicknodeServicesContractAddresses

func (p *ProvisionedQuicknodeServicesContractAddresses) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{p.QuicknodeID, p.ContractAddress}
	}
	return pgValues
}
func (p *ProvisionedQuicknodeServicesContractAddresses) GetTableColumns() (columnValues []string) {
	columnValues = []string{"quicknode_id", "contract_address"}
	return columnValues
}
func (p *ProvisionedQuicknodeServicesContractAddresses) GetTableName() (tableName string) {
	tableName = "provisioned_quicknode_services_contract_addresses"
	return tableName
}
