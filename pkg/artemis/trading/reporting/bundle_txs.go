package artemis_reporting

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
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
	artemis_autogen_bases.EthMempoolMevTx
	TradeExecutionFlow *web3_client.TradeExecutionFlow
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
		rowErr := rows.Scan(&bundle.EthTx.EventID, &bundle.BundleHash, &bundle.EthTx.From, &bundle.EthTx.Nonce, &bundle.EthTx.TxHash,
			&bundle.EthTxGas.GasFeeCap, &bundle.EthTxGas.GasLimit, &bundle.EthTxGas.GasTipCap, &bundle.EthTxGas.GasPrice,
			&bundle.EthTxReceipts.GasUsed, &bundle.EthTxReceipts.EffectiveGasPrice, &bundle.EthTxReceipts.CumulativeGasUsed,
			&bundle.EthTxReceipts.BlockHash, &bundle.EthTxReceipts.TransactionIndex, &bundle.EthTxReceipts.BlockNumber, &bundle.EthTxReceipts.Status,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg("GetBundleSubmissionHistory")
			return bg, rowErr
		}
		bundle.EthTxGas.TxHash = bundle.EthTx.TxHash
		if _, ok := bg.Map[bundle.BundleHash]; !ok {
			bg.Map[bundle.BundleHash] = []Bundle{}
		}
		tmp := bg.Map[bundle.BundleHash]
		tmp = append(tmp, bundle)
		bg.Map[bundle.BundleHash] = tmp
	}
	return bg, nil
}

func getBundlesProfitQ() string {
	var que = `
			WITH cte_bundles AS (
				SELECT eb.event_id, eb.bundle_hash, et."from", et.nonce, et.tx_hash,
					   eg.gas_fee_cap, eg.gas_limit, eg.gas_tip_cap, eg.gas_price,
 					   er.gas_used, er.effective_gas_price, er.cumulative_gas_used,
					   er.block_hash, er.transaction_index, er.block_number, er.status,
				 	   ebp.revenue, ebp.revenue_prediction, ebp.revenue_prediction_skew, ebp.costs, ebp.profit,
					   em.tx_flow_prediction, em.block_number AS block_number_seen
				FROM eth_mev_bundle eb
				INNER JOIN eth_tx et ON et.event_id = eb.event_id
				INNER JOIN eth_tx_gas eg ON eg.tx_hash = et.tx_hash
				INNER JOIN eth_tx_receipts er ON er.tx_hash = et.tx_hash
				LEFT JOIN eth_mev_bundle_profit ebp ON ebp.bundle_hash = eb.bundle_hash
				INNER JOIN eth_mempool_mev_tx em ON em.tx_hash = et.tx_hash
				WHERE eb.event_id > $1 AND eb.protocol_network_id = $2 AND costs > 0
				ORDER BY eb.event_id DESC
			) 
			SELECT *
 			FROM cte_bundles
			`
	return que
}

const (
	AccountAddr = "0x000000641e80A183c8B736141cbE313E136bc8c6"
)

func GetBundlesProfitHistory(ctx context.Context, eventID, protocolNetworkID int) (BundlesGroup, error) {
	bg := BundlesGroup{
		Map: make(map[string][]Bundle),
	}
	q := getBundlesProfitQ()
	rows, err := apps.Pg.Query(ctx, q, eventID, protocolNetworkID)
	if err != nil {
		return bg, err
	}
	defer rows.Close()
	for rows.Next() {
		bundle := Bundle{}
		rowErr := rows.Scan(&bundle.EthTx.EventID, &bundle.BundleHash, &bundle.EthTx.From, &bundle.EthTx.Nonce, &bundle.EthTx.TxHash,
			&bundle.EthTxGas.GasFeeCap, &bundle.EthTxGas.GasLimit, &bundle.EthTxGas.GasTipCap, &bundle.EthTxGas.GasPrice,
			&bundle.EthTxReceipts.GasUsed, &bundle.EthTxReceipts.EffectiveGasPrice, &bundle.EthTxReceipts.CumulativeGasUsed,
			&bundle.EthTxReceipts.BlockHash, &bundle.EthTxReceipts.TransactionIndex, &bundle.EthTxReceipts.BlockNumber, &bundle.EthTxReceipts.Status,
			&bundle.Revenue, &bundle.RevenuePrediction, &bundle.RevenuePrediction, &bundle.Costs, &bundle.Profit,
			&bundle.EthMempoolMevTx.TxFlowPrediction, &bundle.EthMempoolMevTx.BlockNumber,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg("GetBundlesProfitHistory")
			return bg, rowErr
		}
		tfPrediction, terr := web3_client.UnmarshalTradeExecutionFlowToInt(bundle.EthMempoolMevTx.TxFlowPrediction)
		if terr != nil {
			log.Err(terr).Msg("GetBundlesProfitHistory")
			return bg, terr
		}
		bundle.TradeExecutionFlow = tfPrediction
		bundle.EthTxGas.TxHash = bundle.EthTx.TxHash
		if _, ok := bg.Map[bundle.BundleHash]; !ok {
			bg.Map[bundle.BundleHash] = []Bundle{}
		}
		tmp := bg.Map[bundle.BundleHash]
		tmp = append(tmp, bundle)
		bg.Map[bundle.BundleHash] = tmp
	}
	return bg, nil
}
