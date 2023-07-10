package artemis_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type Permit2Tx struct {
	Nonce             int    `db:"nonce" json:"nonce"`
	Owner             string `db:"owner" json:"owner"`
	Deadline          int    `db:"deadline" json:"deadline"`
	EventID           int    `db:"event_id" json:"eventID"`
	Token             string `db:"token" json:"token"`
	ProtocolNetworkID int    `db:"protocol_network_id" json:"protocolNetworkID"`
}
type Permit2TxSlice []Permit2Tx

func (p *Permit2Tx) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{p.Nonce, p.Owner, p.Deadline, p.EventID, p.Token, p.ProtocolNetworkID}
	}
	return pgValues
}
func (p *Permit2Tx) GetTableColumns() (columnValues []string) {
	columnValues = []string{"nonce", "owner", "deadline", "event_id", "token", "protocol_network_id"}
	return columnValues
}
func (p *Permit2Tx) GetTableName() (tableName string) {
	tableName = "permit2_tx"
	return tableName
}
