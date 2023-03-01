package artemis_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type EthScheduledDelivery struct {
	DeliveryScheduleType string `db:"delivery_schedule_type" json:"deliveryScheduleType"`
	ProtocolNetworkID    int    `db:"protocol_network_id" json:"protocolNetworkID"`
	Amount               int    `db:"amount" json:"amount"`
	Units                string `db:"units" json:"units"`
	DeliveryID           int    `db:"delivery_id" json:"deliveryID"`
	PublicKey            string `db:"public_key" json:"publicKey"`
}
type EthScheduledDeliverySlice []EthScheduledDelivery

func (e *EthScheduledDelivery) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{e.DeliveryScheduleType, e.ProtocolNetworkID, e.Amount, e.Units, e.DeliveryID, e.PublicKey}
	}
	return pgValues
}
func (e *EthScheduledDelivery) GetTableColumns() (columnValues []string) {
	columnValues = []string{"delivery_schedule_type", "protocol_network_id", "amount", "units", "delivery_id", "public_key"}
	return columnValues
}
func (e *EthScheduledDelivery) GetTableName() (tableName string) {
	tableName = "eth_scheduled_delivery"
	return tableName
}
