package artemis_eth_txs

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

type EthTx struct {
	artemis_autogen_bases.EthTx
	artemis_autogen_bases.EthTxGas
	Permit2Tx
	NextUserNonce int `json:"nextUserNonce,omitempty"`
}

type Permit2Tx struct {
	artemis_autogen_bases.Permit2Tx
	NextPermit2Nonce int `json:"nextNonce,omitempty"`
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
                    INSERT INTO permit2_tx(event_id, nonce, owner, deadline, "token", protocol_network_id) 
                    VALUES ($1, $11, $12, $13, $14, $3);`
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

func InsertTxsWithBundle(ctx context.Context, dbTx pgx.Tx, txs []EthTx, bundleHash string) error {
	q0 := sql_query_templates.QueryParams{}
	q0.RawQuery = `INSERT INTO events(event_id) VALUES ($1)`
	q1 := sql_query_templates.QueryParams{}
	q1.RawQuery = `WITH cte_tx_1 AS (
                        INSERT INTO eth_tx(event_id, tx_hash, protocol_network_id, nonce, "from", type) 
                        VALUES ($1, $2, $3, $4, $5, $6)
                    ), cte_gas AS (
                        INSERT INTO eth_tx_gas(tx_hash, gas_price, gas_limit, gas_tip_cap, gas_fee_cap)
                        VALUES ($2, $7, $8, $9, $10)
                    ) INSERT INTO permit2_tx(event_id, nonce, owner, deadline, "token", protocol_network_id) 
 					  VALUES ($1, $11, $12, $13, $14, $3)`
	q2 := sql_query_templates.QueryParams{}
	q2.RawQuery = `WITH cte_tx_2 AS (
                        INSERT INTO eth_tx(event_id, tx_hash, protocol_network_id, nonce, "from", type) 
                        VALUES ($1, $2, $3, $4, $5, $6)
                    ) INSERT INTO eth_tx_gas(tx_hash, gas_price, gas_limit, gas_tip_cap, gas_fee_cap)
                      VALUES ($2, $7, $8, $9, $10)`
	q3 := sql_query_templates.QueryParams{}
	q3.RawQuery = `INSERT INTO eth_mev_bundle(event_id, bundle_hash, protocol_network_id)
 			       VALUES ($1, $2, $3)`
	eventID := ts.UnixTimeStampNow()
	for i, e := range txs {
		e.EventID = eventID
		if i == 0 {
			_, err := dbTx.Exec(ctx, q0.RawQuery, e.EventID)
			if err != nil {
				return err
			}
		}
		if e.ProtocolNetworkID == 0 {
			e.ProtocolNetworkID = hestia_req_types.EthereumMainnetProtocolNetworkID
		}
		if e.Type == "" {
			e.Type = "0x02"
		}
		if e.Permit2Tx.Owner == "" || e.Permit2Tx.Token == "" {
			_, err := dbTx.Exec(ctx, q2.RawQuery, e.EventID, e.EthTx.TxHash, e.EthTx.ProtocolNetworkID, e.EthTx.Nonce, e.EthTx.From, e.EthTx.Type,
				e.EthTxGas.GasPrice.Int64, e.EthTxGas.GasLimit.Int64, e.EthTxGas.GasTipCap.Int64, e.EthTxGas.GasFeeCap.Int64)
			if err != nil {
				return err
			}
		} else {
			_, err := dbTx.Exec(ctx, q1.RawQuery, e.EventID, e.EthTx.TxHash, e.EthTx.ProtocolNetworkID, e.EthTx.Nonce, e.EthTx.From, e.EthTx.Type,
				e.EthTxGas.GasPrice.Int64, e.EthTxGas.GasLimit.Int64, e.EthTxGas.GasTipCap.Int64, e.EthTxGas.GasFeeCap.Int64,
				e.Permit2Tx.Nonce, e.Permit2Tx.Owner, e.Permit2Tx.Deadline, e.Permit2Tx.Token)
			if err != nil {
				return err
			}
		}
		if i == len(txs)-1 {
			_, err := dbTx.Exec(ctx, q3.RawQuery, e.EventID, bundleHash, e.EthTx.ProtocolNetworkID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *EthTx) SelectNextUserTxNonce(ctx context.Context) (err error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT COALESCE (MAX(nonce), 0) FROM eth_tx WHERE "from" = $1 AND protocol_network_id = $2;`
	log.Debug().Interface("SelectNextPermit2Nonce", q.LogHeader(ArtemisScheduledDelivery))
	if e.ProtocolNetworkID == 0 {
		e.ProtocolNetworkID = hestia_req_types.EthereumMainnetProtocolNetworkID
	}
	if e.Type == "" {
		e.Type = "0x02"
	}
	err = apps.Pg.QueryRow2(ctx, q.RawQuery, e.From, e.ProtocolNetworkID).Scan(&e.EthTx.Nonce)
	if err != nil {
		return err
	}
	e.NextUserNonce = e.EthTx.Nonce + 1
	return misc.ReturnIfErr(err, q.LogHeader(ArtemisScheduledDelivery))
}

func (pt *Permit2Tx) SelectNextPermit2Nonce(ctx context.Context) (err error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT COALESCE (MAX(nonce), 0) FROM permit2_tx WHERE owner = $1 AND token = $2 AND protocol_network_id = $3;`
	log.Debug().Interface("SelectNextPermit2Nonce", q.LogHeader(ArtemisScheduledDelivery))
	if pt.ProtocolNetworkID == 0 {
		pt.ProtocolNetworkID = hestia_req_types.EthereumMainnetProtocolNetworkID
	}
	err = apps.Pg.QueryRow2(ctx, q.RawQuery, pt.Owner, pt.Token, pt.ProtocolNetworkID).Scan(&pt.Nonce)
	if err != nil {
		return err
	}
	pt.NextPermit2Nonce = pt.Nonce + 1
	return misc.ReturnIfErr(err, q.LogHeader(ArtemisScheduledDelivery))
}
