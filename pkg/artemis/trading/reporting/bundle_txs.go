package artemis_reporting

import (
	"context"
	"fmt"

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
				ORDER BY eb.event_id DESC, et."from", et.nonce ASC
			) 
			SELECT *
 			FROM cte_bundles
			`
	return que
}

type BundlesGroup struct {
	Map             map[string][]Bundle `json:"bundles"`
	bundleHashOrder []string
	bundleHashToId  map[string]int
}

type Bundle struct {
	artemis_autogen_bases.EthMevBundleProfit `json:"ethMevBundleProfit,omitempty"`
	artemis_autogen_bases.EthTx              `json:"ethTx,omitempty"`
	artemis_autogen_bases.EthTxGas           `json:"ethTxGas,omitempty"`
	artemis_autogen_bases.EthTxReceipts      `json:"ethTxReceipts,omitempty"`
	artemis_autogen_bases.EthMempoolMevTx    `json:"ethMempoolMevTx,omitempty"`
	TradeExecutionFlow                       *web3_client.TradeExecutionFlow `json:"tradeExecutionFlow"`
}

func (b *Bundle) PrintBundleInfo() {
	fmt.Println("===============================================================================================================")
	fmt.Println("===============================================================================================================")
	fmt.Println("bundleTx.EthTx.EventID", b.EthTx.EventID)
	fmt.Println("bundleTx.BundleHash", b.BundleHash)
	fmt.Println("===============================================================================================================")
	fmt.Println("bundleTx.EthTx.TxHash", b.EthTx.TxHash)
	fmt.Println("bundleTx.EthTxGas.GasTipCap", b.EthTxGas.GasTipCap)
	fmt.Println("bundleTx.EthTxGas.GasFeeCap", b.EthTxGas.GasFeeCap)
	fmt.Println("bundleTx.EthTxGas.GasLimit", b.EthTxGas.GasLimit)
	fmt.Println("bundleTx.EthTxReceipts.EffectiveGasPrice", b.EffectiveGasPrice)
	fmt.Println("bundleTx.EthTxReceipts.BlockNumberRx", b.EthTxReceipts.BlockNumber)

	if b.TradeExecutionFlow != nil {
		fmt.Println("bundleTx.TradeExecutionFlow.CurrentBlockNumber", b.TradeExecutionFlow.CurrentBlockNumber)
		fmt.Println("bundleTx.TradeExecutionFlow.Trade.TradeMethod", b.TradeExecutionFlow.Trade.TradeMethod)
		fmt.Println("bundleTx.TradeExecutionFlow.FrontRunTrade.AmountIn", b.TradeExecutionFlow.FrontRunTrade.AmountIn.String())
		fmt.Println("bundleTx.TradeExecutionFlow.SandwichPrediction.SellAmount", b.TradeExecutionFlow.SandwichPrediction.SellAmount)
		fmt.Println("bundleTx.TradeExecutionFlow.SandwichPrediction.ExpectedProfit", b.TradeExecutionFlow.SandwichPrediction.ExpectedProfit)
		fmt.Println("bundleTx.TradeExecutionFlow.SandwichTrade.AmountOut", b.TradeExecutionFlow.SandwichTrade.AmountOut.String())
	}
	fmt.Println("bundleTx.EthMevBundleProfit.Revenue", b.Revenue)
	fmt.Println("bundleTx.EthMevBundleProfit.RevenuePrediction", b.RevenuePrediction)
	fmt.Println("bundleTx.EthMevBundleProfit.RevenuePredictionSkew", b.RevenuePredictionSkew)
	fmt.Println("bundleTx.EthMevBundleProfit.Costs", b.Costs)
	fmt.Println("bundleTx.EthMevBundleProfit.Profit", b.Profit)
	fmt.Println("===============================================================================================================")
	fmt.Println("===============================================================================================================")
}

func GetBundleSubmissionHistory(ctx context.Context, eventID, protocolNetworkID int) (BundlesGroup, error) {
	bg := BundlesGroup{
		Map:             make(map[string][]Bundle),
		bundleHashOrder: make([]string, 0),
		bundleHashToId:  make(map[string]int),
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

		if _, ok := bg.bundleHashToId[bundle.BundleHash]; !ok {
			bg.bundleHashToId[bundle.BundleHash] = bundle.EthTx.EventID
			bg.bundleHashOrder = append(bg.bundleHashOrder, bundle.BundleHash)
		} else {
			continue
		}
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
			&bundle.EthMevBundleProfit.Revenue, &bundle.EthMevBundleProfit.RevenuePrediction, &bundle.EthMevBundleProfit.RevenuePredictionSkew, &bundle.Costs, &bundle.Profit,
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
