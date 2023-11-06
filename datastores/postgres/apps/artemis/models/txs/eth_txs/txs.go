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
		e.EthTxGas.GasPrice, e.EthTxGas.GasLimit, e.EthTxGas.GasTipCap, e.EthTxGas.GasFeeCap,
		pt.Nonce, pt.Owner, pt.Deadline, pt.Token)
	if err != nil {
		log.Err(err).Msg("InsertTx")
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader(ArtemisScheduledDelivery))
}

func InsertTxsWithBundle(ctx context.Context, dbTx pgx.Tx, txs []EthTx, bundleHash string) error {
	q0 := sql_query_templates.QueryParams{}
	q0.RawQuery = `INSERT INTO events(event_id) VALUES ($1)`
	q1 := sql_query_templates.QueryParams{}
	q1.RawQuery = `WITH cte_tx_1 AS (
						INSERT INTO eth_tx(event_id, tx_hash, protocol_network_id, nonce, "from", type, nonce_id) 
                        SELECT $1, $2, $3, $4, $5, $6, $15
						ON CONFLICT (tx_hash)
						DO UPDATE SET	
					 		nonce_id = EXCLUDED.nonce_id
                        RETURNING *
                    ), cte_gas AS (
                        INSERT INTO eth_tx_gas(tx_hash, gas_price, gas_limit, gas_tip_cap, gas_fee_cap)
                        VALUES ($2, $7, $8, $9, $10)
						ON CONFLICT (tx_hash) 
						DO UPDATE SET 
							gas_price = EXCLUDED.gas_price,
							gas_limit = EXCLUDED.gas_limit,
							gas_tip_cap = EXCLUDED.gas_tip_cap,
							gas_fee_cap = EXCLUDED.gas_fee_cap
                    ) INSERT INTO permit2_tx(event_id, nonce, owner, deadline, "token", protocol_network_id, nonce_id) 
 					  VALUES ($1, $11, $12, $13, $14, $3, $15);`
	q2 := sql_query_templates.QueryParams{}
	q2.RawQuery = `WITH cte_tx_2 AS (
						INSERT INTO eth_tx(event_id, tx_hash, protocol_network_id, nonce, "from", type, nonce_id) 
                        SELECT $1, $2, $3, $4, $5, $6, $11
						ON CONFLICT (tx_hash)
						DO UPDATE SET	
					 		nonce_id = EXCLUDED.nonce_id
                        RETURNING *
                    ) 	
						INSERT INTO eth_tx_gas(tx_hash, gas_price, gas_limit, gas_tip_cap, gas_fee_cap)
						VALUES ($2, $7, $8, $9, $10)
						ON CONFLICT (tx_hash) 
						DO UPDATE SET 
							gas_price = EXCLUDED.gas_price,
							gas_limit = EXCLUDED.gas_limit,
							gas_tip_cap = EXCLUDED.gas_tip_cap,
							gas_fee_cap = EXCLUDED.gas_fee_cap;
`
	q3 := sql_query_templates.QueryParams{}
	q3.RawQuery = `INSERT INTO eth_mev_bundle(event_id, bundle_hash, protocol_network_id)
 			       VALUES ($1, $2, $3)`
	eventID := ts.UnixTimeStampNow()
	for i, e := range txs {
		e.EventID = eventID
		if i == 0 {
			_, err := dbTx.Exec(ctx, q0.RawQuery, e.EventID)
			if err != nil {
				log.Err(err).Msg("InsertTxsWithBundle")
				return err
			}
		}
		if e.ProtocolNetworkID == 0 {
			e.ProtocolNetworkID = hestia_req_types.EthereumMainnetProtocolNetworkID
		}
		if e.Type == "" {
			e.Type = "0x02"
		}
		// query 2
		if e.Permit2Tx.Owner == "" || e.Permit2Tx.Token == "" {
			_, err := dbTx.Exec(ctx, q2.RawQuery, e.EventID, e.EthTx.TxHash, e.EthTx.ProtocolNetworkID, e.EthTx.Nonce, e.EthTx.From, e.EthTx.Type,
				e.EthTxGas.GasPrice, e.EthTxGas.GasLimit, e.EthTxGas.GasTipCap, e.EthTxGas.GasFeeCap, ts.UnixTimeStampNow())
			if err != nil {
				log.Warn().Interface("tx", e).Msg("InsertTxsWithBundle: Query2")
				log.Err(err).Msg("InsertTxsWithBundle: Query2")
				return err
			}
		} else {
			// query 1
			_, err := dbTx.Exec(ctx, q1.RawQuery, e.EventID, e.EthTx.TxHash, e.EthTx.ProtocolNetworkID, e.EthTx.Nonce, e.EthTx.From, e.EthTx.Type,
				e.EthTxGas.GasPrice, e.EthTxGas.GasLimit, e.EthTxGas.GasTipCap, e.EthTxGas.GasFeeCap,
				e.Permit2Tx.Nonce, e.Permit2Tx.Owner, e.Permit2Tx.Deadline, e.Permit2Tx.Token, ts.UnixTimeStampNow())
			if err != nil {
				log.Warn().Interface("tx", e).Msg("InsertTxsWithBundle: Query1")
				log.Err(err).Msg("InsertTxsWithBundle")
				return err
			}
		}
		if i == len(txs)-1 {
			_, err := dbTx.Exec(ctx, q3.RawQuery, e.EventID, bundleHash, e.EthTx.ProtocolNetworkID)
			if err != nil {
				log.Err(err).Msg("InsertTxsWithBundle")
				return err
			}
		}
	}
	return nil
}

func (e *EthTx) SelectNextUserTxNonce(ctx context.Context) (err error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT COALESCE (MAX(nonce_id), 0), nonce FROM eth_tx WHERE "from" = $1 AND protocol_network_id = $2 GROUP BY nonce;`
	log.Debug().Interface("SelectNextUserTxNonce", q.LogHeader("SelectNextUserTxNonce"))
	if e.ProtocolNetworkID == 0 {
		e.ProtocolNetworkID = hestia_req_types.EthereumMainnetProtocolNetworkID
	}
	if e.Type == "" {
		e.Type = "0x02"
	}
	tmp := 0
	err = apps.Pg.QueryRow2(ctx, q.RawQuery, e.From, e.ProtocolNetworkID).Scan(&tmp, &e.EthTx.Nonce)
	if err == pgx.ErrNoRows {
		e.Nonce = 0
		err = nil
	}
	if err != nil {
		return err
	}
	e.NextUserNonce = e.EthTx.Nonce + 1
	return misc.ReturnIfErr(err, q.LogHeader("ArtemisScheduledDelivery"))
}

func (pt *Permit2Tx) SelectNextPermit2Nonce(ctx context.Context) (err error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT COALESCE (MAX(nonce_id), 0), nonce FROM permit2_tx WHERE owner = $1 AND token = $2 AND protocol_network_id = $3 GROUP BY nonce;`
	log.Debug().Interface("SelectNextPermit2Nonce", q.LogHeader("SelectNextPermit2Nonce"))
	if pt.ProtocolNetworkID == 0 {
		pt.ProtocolNetworkID = hestia_req_types.EthereumMainnetProtocolNetworkID
	}
	tmp := 0
	err = apps.Pg.QueryRow2(ctx, q.RawQuery, pt.Owner, pt.Token, pt.ProtocolNetworkID).Scan(&tmp, &pt.Nonce)
	if err == pgx.ErrNoRows {
		pt.Nonce = 0
		err = nil
		return err
	}
	if err != nil {
		return err
	}
	pt.NextPermit2Nonce = pt.Nonce + 1
	return misc.ReturnIfErr(err, q.LogHeader("SelectNextPermit2Nonce"))
}

type TxRxEvent struct {
	EventID int
	TxHash  string
}

func SelectExternalTxs(ctx context.Context, traderAddress string, protocolID, minEventID int) ([]TxRxEvent, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT et.event_id, et.tx_hash
					FROM eth_tx et
					WHERE NOT EXISTS (
					  SELECT 1
					  FROM eth_tx_receipts er
					  WHERE er.tx_hash = et.tx_hash
					)
					AND et."from" != $1 AND et.protocol_network_id = $2 AND et.event_id >= $3
					AND NOT EXISTS (
					  SELECT 1
					  FROM eth_tx_receipts er
					  WHERE er.tx_hash = et.tx_hash AND (er.status = 'failed' OR er.status = 'success')
					)
					ORDER BY et.event_id DESC;`

	rows, err := apps.Pg.Query(ctx, q.RawQuery, traderAddress, protocolID, minEventID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("SelectExternalTxs")); returnErr != nil {
		return nil, err
	}
	defer rows.Close()
	var rxEvents []TxRxEvent
	for rows.Next() {
		rxEvent := TxRxEvent{}
		rowErr := rows.Scan(
			&rxEvent.EventID, &rxEvent.TxHash,
		)
		if rowErr != nil {
			return nil, rowErr
		}
		rxEvents = append(rxEvents, rxEvent)
	}
	return rxEvents, nil
}
