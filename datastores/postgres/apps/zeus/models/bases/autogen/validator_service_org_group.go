package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ValidatorServiceOrgGroup struct {
	GroupName         string `db:"group_name" json:"groupName"`
	OrgID             int    `db:"org_id" json:"orgID"`
	Pubkey            string `db:"pubkey" json:"pubkey"`
	ProtocolNetworkID int    `db:"protocol_network_id" json:"protocolNetworkID"`
	FeeRecipient      string `db:"fee_recipient" json:"feeRecipient"`
}
type ValidatorServiceOrgGroupSlice []ValidatorServiceOrgGroup

func (v *ValidatorServiceOrgGroup) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{v.GroupName, v.OrgID, v.Pubkey, v.ProtocolNetworkID, v.FeeRecipient}
	}
	return pgValues
}
func (v *ValidatorServiceOrgGroup) GetTableColumns() (columnValues []string) {
	columnValues = []string{"group_name", "org_id", "pubkey", "protocol_network_id", "fee_recipient"}
	return columnValues
}
func (v *ValidatorServiceOrgGroup) GetTableName() (tableName string) {
	tableName = "validator_service_org_group"
	return tableName
}
