package hestia_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ValidatorServiceOrgGroup struct {
	GroupName         string `db:"group_name" json:"groupName"`
	OrgID             *int   `db:"org_id" json:"orgID,omitempty"`
	Pubkey            string `db:"pubkey" json:"pubkey"`
	ProtocolNetworkID int    `db:"protocol_network_id" json:"protocolNetworkID"`
	FeeRecipient      string `db:"fee_recipient" json:"feeRecipient"`
	Enabled           bool   `db:"enabled" json:"enabled"`
	ServiceURL        string `db:"service_url" json:"serviceURL"`
}

type ValidatorServiceOrgGroupSlice []ValidatorServiceOrgGroup

func (v *ValidatorServiceOrgGroup) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{v.GroupName, v.OrgID, v.Pubkey, v.ProtocolNetworkID, v.FeeRecipient, v.Enabled, v.ServiceURL}
	}
	return pgValues
}
func (v *ValidatorServiceOrgGroup) GetTableColumns() (columnValues []string) {
	columnValues = []string{"group_name", "org_id", "pubkey", "protocol_network_id", "fee_recipient", "enabled", "service_url"}
	return columnValues
}

func (v *ValidatorServiceOrgGroup) GetTableName() (tableName string) {
	tableName = "validators_service_org_groups"
	return tableName
}
