package artemis_autogen_bases

import (
	"database/sql"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

type EthEcdsaTxs struct {
	TxID                 int            `db:"tx_id" json:"txID"`
	ProtocolNetworkID    int            `db:"protocol_network_id" json:"protocolNetworkID"`
	PublicKeyTypeID      int            `db:"public_key_type_id" json:"publicKeyTypeID"`
	PublicKey            string         `db:"public_key" json:"publicKey"`
	Nonce                int            `db:"nonce" json:"nonce"`
	AmountGwei           sql.NullInt64  `db:"amount_gwei" json:"amountGwei"`
	AmountGweiDecimals   sql.NullInt64  `db:"amount_gwei_decimals" json:"amountGweiDecimals"`
	GasLimitGwei         sql.NullInt64  `db:"gas_limit_gwei" json:"gasLimitGwei"`
	GasLimitGweIDecimals sql.NullInt64  `db:"gas_limit_gwei_decimals" json:"gasLimitGweiDecimals"`
	GasPriceGwei         sql.NullInt64  `db:"gas_price_gwei" json:"gasPriceGwei"`
	GasPriceGweiDecimals sql.NullInt64  `db:"gas_price_gwei_decimals" json:"gasPriceGweiDecimals"`
	Payload              sql.NullString `db:"payload" json:"payload"`
	To                   sql.NullString `db:"to" json:"to"`
	TxHash               sql.NullString `db:"tx_hash" json:"txHash"`
	R                    string         `db:"r" json:"r"`
	S                    string         `db:"s" json:"s"`
	V                    string         `db:"v" json:"v"`
}
type EthEcdsaTxsSlice []EthEcdsaTxs

func (e *EthEcdsaTxs) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{e.GasPriceGwei, e.GasLimitGweIDecimals, e.R, e.TxHash, e.TxID, e.AmountGwei, e.PublicKeyTypeID, e.PublicKey, e.GasLimitGwei, e.S, e.Payload, e.To, e.ProtocolNetworkID, e.Nonce, e.GasPriceGweiDecimals, e.AmountGweiDecimals, e.V}
	}
	return pgValues
}
func (e *EthEcdsaTxs) GetTableColumns() (columnValues []string) {
	columnValues = []string{"gas_price_gwei", "gas_limit_gwei_decimals", "r", "tx_hash", "tx_id", "amount_gwei", "public_key_type_id", "public_key", "gas_limit_gwei", "s", "payload", "to", "protocol_network_id", "nonce", "gas_price_gwei_decimals", "amount_gwei_decimals", "v"}
	return columnValues
}
func (e *EthEcdsaTxs) GetTableName() (tableName string) {
	tableName = "eth_ecdsa_txs"
	return tableName
}
