package artemis_reporting

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
)

func getBundlesQ() string {
	var que = `
			WITH cte_bundles AS (
				SELECT eb.event_id, eb.bundle_hash, et."from", et.nonce, et.tx_hash,
					   eg.gas_fee_cap, eg.gas_limit, eg.gas_tip_cap, eg.gas_price,
 					   er.gas_used, er.effective_gas_price, er.cumulative_gas_used,
					   er.block_hash, er.transaction_index, er.block_number, er.status
				FROM eth_mev_bundle eb
				INNER JOIN eth_tx et ON et.event_id = eb.event_id
				INNER JOIN eth_tx_gas eg ON eg.tx_hash = et.tx_hash
				INNER JOIN eth_tx_receipts er ON er.tx_hash = et.tx_hash
				WHERE eb.event_id > $1 AND eb.protocol_network_id = $2 
				ORDER BY eb.event_id DESC, et.nonce ASC
			) 
			SELECT *
 			FROM cte_bundles
			`
	return que
}

type BundlesGroup struct {
	Map map[string][]Bundle
}

type Bundle struct {
	artemis_autogen_bases.EthMevBundleProfit
	artemis_autogen_bases.EthTx
	artemis_autogen_bases.EthTxGas
	artemis_autogen_bases.EthTxReceipts
}

func GetBundleSubmissionHistory(ctx context.Context, eventID, protocolNetworkID int) (BundlesGroup, error) {
	bg := BundlesGroup{
		Map: make(map[string][]Bundle),
	}
	q := getBundlesQ()
	rows, err := apps.Pg.Query(ctx, q, eventID, protocolNetworkID)
	if err != nil {
		return bg, err
	}
	defer rows.Close()
	for rows.Next() {
		bundle := Bundle{}
		bundleHash := ""
		rowErr := rows.Scan(&bundle.EthTx.EventID, &bundleHash, &bundle.From, &bundle.Nonce, &bundle.EthTx.TxHash,
			&bundle.EthTxGas.GasFeeCap, &bundle.EthTxGas.GasLimit, &bundle.EthTxGas.GasTipCap, &bundle.EthTxGas.GasPrice,
			&bundle.EthTxReceipts.GasUsed, &bundle.EthTxReceipts.EffectiveGasPrice, &bundle.EthTxReceipts.CumulativeGasUsed,
			&bundle.EthTxReceipts.BlockHash, &bundle.EthTxReceipts.TransactionIndex, &bundle.EthTxReceipts.BlockNumber, &bundle.EthTxReceipts.Status,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg("GetBundleSubmissionHistory")
			return bg, rowErr
		}
		bundle.EthTxGas.TxHash = bundle.EthTx.TxHash
		if _, ok := bg.Map[bundleHash]; !ok {
			bg.Map[bundleHash] = []Bundle{}
		}
		tmp := bg.Map[bundleHash]
		tmp = append(tmp, bundle)
		bg.Map[bundleHash] = tmp
	}
	return bg, nil
}
