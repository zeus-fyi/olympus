package artemis_eth_txs

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

/*
type EthTx struct {
	ProtocolNetworkID int    `db:"protocol_network_id" json:"protocolNetworkID"`
	TxHash            string `db:"tx_hash" json:"txHash"`
	Nonce             int    `db:"nonce" json:"nonce"`
	From              string `db:"from" json:"from"`
	Type              string `db:"type" json:"type"`
	EventID           int    `db:"event_id" json:"eventID"`
}

type EthTxGas struct {
	TxHash    string        `db:"tx_hash" json:"txHash"`
	GasPrice  sql.NullInt64 `db:"gasPrice" json:"gasPrice"`
	GasLimit  sql.NullInt64 `db:"gasLimit" json:"gasLimit"`
	GasTipCap sql.NullInt64 `db:"gasTipCap" json:"gasTipCap"`
	GasFeeCap sql.NullInt64 `db:"gasFeeCap" json:"gasFeeCap"`
}


type Permit2Tx struct {
	Nonce    int    `db:"nonce" json:"nonce"`
	Owner    string `db:"owner" json:"owner"`
	Deadline int    `db:"deadline" json:"deadline"`
	EventID  int    `db:"event_id" json:"eventID"`
	Token    string `db:"token" json:"token"`
}
*/

type EthTx struct {
	artemis_autogen_bases.EthTx
	artemis_autogen_bases.EthTxGas
	NextNonce int `db:"next_nonce" json:"nextNonce"`
}

type Permit2Tx struct {
	artemis_autogen_bases.Permit2Tx
}

const ArtemisScheduledDelivery = "EthTx"

var ts chronos.Chronos

func (e *EthTx) InsertTx(ctx context.Context, pt Permit2Tx) (err error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `WITH cte_insert_id AS (
                        INSERT INTO events(event_id) VALUES ($1) RETURNING event_id
                    ), cte_tx AS (
                        INSERT INTO eth_tx(event_id, tx_hash, protocol_network_id, nonce, "from", type) 
                        VALUES ($1, $2, $3, $4, $5, $6) RETURNING event_id
                    ), cte_gas AS (
                        INSERT INTO eth_tx_gas(tx_hash, gas_price, gas_limit, gas_tip_cap, gas_fee_cap)
                        VALUES ($2, $7, $8, $9, $10) RETURNING tx_hash
                    )
                    INSERT INTO permit2_tx(event_id, nonce, owner, deadline, "token") 
                    VALUES ($1, $11, $12, $13, $14);`
	log.Debug().Interface("InsertTx", q.LogHeader(ArtemisScheduledDelivery))
	if e.ProtocolNetworkID == 0 {
		e.ProtocolNetworkID = hestia_req_types.EthereumMainnetProtocolNetworkID
	}
	if e.Type == "" {
		e.Type = "0x02"
	}
	e.EventID = ts.UnixTimeStampNow()
	_, err = apps.Pg.Exec(ctx, q.RawQuery, e.EventID, e.EthTx.TxHash, e.EthTx.ProtocolNetworkID, e.EthTx.Nonce, e.EthTx.From, e.EthTx.Type,
		e.EthTxGas.GasPrice.Int64, e.EthTxGas.GasLimit.Int64, e.EthTxGas.GasTipCap.Int64, e.EthTxGas.GasFeeCap.Int64,
		pt.Nonce, pt.Owner, pt.Deadline, pt.Token)
	if err != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader(ArtemisScheduledDelivery))
}

func (e *EthTx) SelectNextTxNonce(ctx context.Context, pt Permit2Tx) (err error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT COALESCE (MAX(nonce), 0) FROM permit2_tx WHERE owner = $1 AND token = $2;`
	log.Debug().Interface("SelectNextTxNonce", q.LogHeader(ArtemisScheduledDelivery))
	if e.ProtocolNetworkID == 0 {
		e.ProtocolNetworkID = hestia_req_types.EthereumMainnetProtocolNetworkID
	}
	if e.Type == "" {
		e.Type = "0x02"
	}
	err = apps.Pg.QueryRow2(ctx, q.RawQuery, pt.Owner, pt.Token).Scan(&e.EthTx.Nonce)
	if err != nil {
		return err
	}
	e.NextNonce = e.EthTx.Nonce + 1
	return misc.ReturnIfErr(err, q.LogHeader(ArtemisScheduledDelivery))
}
