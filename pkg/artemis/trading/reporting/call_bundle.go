package artemis_reporting

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"unicode"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
	"github.com/metachris/flashbotsrpc"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func getCallBundleSaveQ() string {
	var que = `
        WITH cte_mev_call AS (
			INSERT INTO events (event_id)
			VALUES ($1)
			RETURNING event_id
		),
		cte_eth_tx AS (
			INSERT INTO eth_tx(event_id, tx_hash, protocol_network_id, nonce, "from", type, nonce_id)
			SELECT 
				c.event_id, 
				$6 as tx_hash,
				$4 as protocol_network_id,
				$7 as nonce,
				$9 as "from", 
				$8 as type,
				$1 as nonce_id
			FROM cte_mev_call c
			RETURNING event_id
		)
		INSERT INTO eth_mev_call_bundle (event_id, builder_name, bundle_hash, protocol_network_id, eth_call_resp_json)
		SELECT 
			c.event_id, 
			$2 as builder_name,
			$3 as bundle_hash,
			$4 as protocol_network_id, 
			$5::jsonb as eth_call_resp_json
		FROM cte_mev_call c;`
	return que
}

func clearString(str string) string {
	str = strings.Replace(str, "�y�", "", -1)
	var result []rune
	for _, r := range str {
		// Keep only letters, digits, and underscore
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			result = append(result, r)
		}
	}
	return string(result)
}

func InsertCallBundleResp(ctx context.Context, builder string, protocolID int, callBundlesResp flashbotsrpc.FlashbotsCallBundleResponse, tf *web3_client.TradeExecutionFlow) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = getCallBundleSaveQ()
	if callBundlesResp.BundleHash == "" {
		return errors.New("bundle hash is empty")
	}
	if tf.Tx == nil {
		return errors.New("tx is nil")
	}

	tmp := flashbotsrpc.FlashbotsCallBundleResponse{
		BundleGasPrice:    callBundlesResp.BundleGasPrice,
		BundleHash:        callBundlesResp.BundleHash,
		CoinbaseDiff:      callBundlesResp.CoinbaseDiff,
		EthSentToCoinbase: callBundlesResp.EthSentToCoinbase,
		GasFees:           callBundlesResp.GasFees,
		Results:           make([]flashbotsrpc.FlashbotsCallBundleResult, len(callBundlesResp.Results)),
		StateBlockNumber:  callBundlesResp.StateBlockNumber,
		TotalGasUsed:      callBundlesResp.TotalGasUsed,
	}

	for i, v := range callBundlesResp.Results {
		fr := flashbotsrpc.FlashbotsCallBundleResult{
			CoinbaseDiff:      v.CoinbaseDiff,
			EthSentToCoinbase: v.EthSentToCoinbase,
			FromAddress:       v.FromAddress,
			GasFees:           v.GasFees,
			GasPrice:          v.GasPrice,
			GasUsed:           v.GasUsed,
			ToAddress:         v.ToAddress,
			TxHash:            v.TxHash,
			Value:             v.Value,
			Error:             clearString(v.Error),
			Revert:            clearString(v.Revert),
		}
		tmp.Results[i] = fr
	}

	b, err := json.Marshal(tmp)
	if err != nil {
		return err
	}
	ts := chronos.Chronos{}
	eventID := ts.UnixTimeStampNow()

	typeTxStr := "0x02"
	typeTx := tf.Tx.Type()
	if typeTx == 1 {
		typeTxStr = "0x01"
	}

	fromStr := ""
	chainId := artemis_eth_units.NewBigInt(hestia_req_types.EthereumMainnetProtocolNetworkID)
	if tf.Tx.ChainId() != nil {
		chainId = tf.Tx.ChainId()
	}
	sender := types.LatestSignerForChainID(chainId)
	from, ferr := sender.Sender(tf.Tx)
	if ferr != nil {
		log.Err(ferr).Msg("failed to get sender")
		return ferr
	} else {
		fromStr = from.String()
	}
	_, err = apps.Pg.Exec(ctx, q.RawQuery, eventID, builder, callBundlesResp.BundleHash, protocolID, string(b), tf.Tx.Hash().String(), tf.Tx.Nonce(), typeTxStr, fromStr)
	if err == pgx.ErrNoRows {
		err = nil
		return err
	} else {
		log.Warn().Interface("tx", tf.Tx).Msg("InsertTxsWithBundle: InsertCallBundleResp")
	}
	return misc.ReturnIfErr(err, q.LogHeader("InsertBundleProfit"))
}

func selectCallBundles() string {
	var que = `SELECT eb.event_id, builder_name, bundle_hash, eth_call_resp_json,
				etx.tx_hash, etx."from",
				mem.tx_flow_prediction, ea.amount_in, ea.trade_method, ea.expected_profit_amount_out, ea.actual_profit_amount_out, 
				er.effective_gas_price, er.gas_used, er.status, er.block_number, er.transaction_index
				FROM eth_mev_call_bundle eb
				INNER JOIN eth_tx etx ON etx.event_id = eb.event_id
				INNER JOIN eth_mev_tx_analysis ea ON ea.tx_hash = etx.tx_hash
				INNER JOIN eth_mempool_mev_tx mem ON mem.tx_hash = etx.tx_hash
				INNER JOIN eth_tx_receipts er ON er.tx_hash = etx.tx_hash
			   WHERE eb.event_id > $1 AND eb.protocol_network_id = $2
			   ORDER BY eb.event_id DESC
			   LIMIT 100;
		  	   	`
	return que
}

type CallBundleHistory struct {
	EventID                                  int                                      `json:"eventID"`
	BuilderName                              string                                   `json:"builderName"`
	BundleHash                               string                                   `json:"bundleHash"`
	TxHash                                   string                                   `json:"txHash"`
	FromAddress                              string                                   `json:"from"`
	TradeExecutionFlowJSON                   web3_client.TradeExecutionFlowJSON       `json:"tradeExecutionFlowJSON"`
	Trades                                   []artemis_trading_types.JSONTradeOutcome `json:"trades"`
	PairAddress                              string                                   `json:"pairAddress"`
	AmountIn                                 string                                   `json:"amountIn"` // Assuming numeric field
	SeenAtBlockNumber                        int                                      `json:"seenAtBlockNumber"`
	TradeMethod                              string                                   `json:"tradeMethod"`
	ExpectedProfitAmountOut                  string                                   `json:"expectedProfitAmountOut"` // Assuming numeric field
	ActualProfitAmountOut                    string                                   `json:"actualProfitAmountOut"`   // Assuming numeric field
	EffectiveGasPrice                        int                                      `json:"effectiveGasPrice"`       // Assuming this is an integer value
	GasUsed                                  int                                      `json:"gasUsed"`
	Status                                   string                                   `json:"status"`
	BlockNumber                              int                                      `json:"blockNumber"`
	TransactionIndex                         int                                      `json:"transactionIndex"`
	SubmissionTime                           string                                   `json:"submissionTime"` // Already present in your struct
	flashbotsrpc.FlashbotsCallBundleResponse `json:"flashbotsCallBundleResponse"`     // Embedded struct, ensure fields are mapped correctly
}

func SelectCallBundleHistory(ctx context.Context, minEventId, protocolNetworkID int) ([]CallBundleHistory, error) {
	rows, err := apps.Pg.Query(ctx, selectCallBundles(), minEventId, protocolNetworkID)
	if err != nil {
		return nil, err
	}
	var rw []CallBundleHistory
	ts := chronos.Chronos{}

	defer rows.Close()
	for rows.Next() {
		cbh := CallBundleHistory{
			FlashbotsCallBundleResponse: flashbotsrpc.FlashbotsCallBundleResponse{},
		}
		txFlow := ""
		rowErr := rows.Scan(
			&cbh.EventID, &cbh.BuilderName, &cbh.BundleHash, &cbh.FlashbotsCallBundleResponse,
			&cbh.TxHash, &cbh.FromAddress,
			&txFlow, &cbh.AmountIn, &cbh.TradeMethod, &cbh.ExpectedProfitAmountOut, &cbh.ActualProfitAmountOut,
			&cbh.EffectiveGasPrice, &cbh.GasUsed, &cbh.Status, &cbh.BlockNumber, &cbh.TransactionIndex,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg("SelectCallBundleHistory")
			return nil, rowErr
		}
		txFlowJson, berr := web3_client.UnmarshalTradeExecutionFlow(txFlow)
		if berr != nil {
			log.Err(berr).Msg("SelectCallBundleHistory")
			return nil, berr
		}
		cbh.TradeExecutionFlowJSON = txFlowJson
		cbh.SeenAtBlockNumber = int(txFlowJson.CurrentBlockNumber.Int64())
		if txFlowJson.InitialPair != nil {
			cbh.PairAddress = txFlowJson.InitialPair.PairContractAddr
		}
		if txFlowJson.InitialPairV3 != nil {
			cbh.PairAddress = txFlowJson.InitialPairV3.PoolAddress
		}

		cbh.Trades = []artemis_trading_types.JSONTradeOutcome{
			txFlowJson.FrontRunTrade,
			txFlowJson.UserTrade,
			txFlowJson.SandwichTrade,
		}

		// Create a big.Float representation of 1e9 for the division to gwei
		weiToGwei := new(big.Float).SetFloat64(1e9)

		// Divide the gas price by 1e18 to convert wei to ether
		bundleGasPriceWei := artemis_eth_units.NewBigFloatFromStr(cbh.FlashbotsCallBundleResponse.BundleGasPrice)
		bundleGasPrice, _ := new(big.Float).Quo(bundleGasPriceWei, weiToGwei).Float64()
		// Format the float to a string with 5 decimal places
		cbh.BundleGasPrice = fmt.Sprintf("%.5f", bundleGasPrice)

		eth := new(big.Float).SetFloat64(1e18)
		bundleGasFeesWei := artemis_eth_units.NewBigFloatFromStr(cbh.FlashbotsCallBundleResponse.GasFees)
		bundleGasFees, _ := new(big.Float).Quo(bundleGasFeesWei, eth).Float64()
		cbh.GasFees = fmt.Sprintf("%.5f", bundleGasFees)

		ethCoinbaseDiffWei := artemis_eth_units.NewBigFloatFromStr(cbh.FlashbotsCallBundleResponse.CoinbaseDiff)
		ethCoinbaseDiff, _ := new(big.Float).Quo(ethCoinbaseDiffWei, eth).Float64()
		cbh.CoinbaseDiff = fmt.Sprintf("%.5f", ethCoinbaseDiff)

		expProfitWei := artemis_eth_units.NewBigFloatFromStr(cbh.ExpectedProfitAmountOut)
		expProfit, _ := new(big.Float).Quo(expProfitWei, eth).Float64()
		cbh.ExpectedProfitAmountOut = fmt.Sprintf("%.5f", expProfit)

		actualProfitWei := artemis_eth_units.NewBigFloatFromStr(cbh.ActualProfitAmountOut)
		actualProfit, _ := new(big.Float).Quo(actualProfitWei, eth).Float64()
		cbh.ActualProfitAmountOut = fmt.Sprintf("%.5f", actualProfit)

		for i, v := range cbh.FlashbotsCallBundleResponse.Results {
			bundleGasPriceWei = artemis_eth_units.NewBigFloatFromStr(v.GasPrice)
			bundleGasPrice, _ = new(big.Float).Quo(bundleGasPriceWei, weiToGwei).Float64()
			// Format the float to a string with 5 decimal places
			cbh.FlashbotsCallBundleResponse.Results[i].GasPrice = fmt.Sprintf("%.5f Gwei", bundleGasPrice)

			ethCoinbaseDiffWei = artemis_eth_units.NewBigFloatFromStr(v.CoinbaseDiff)
			ethCoinbaseDiff, _ = new(big.Float).Quo(ethCoinbaseDiffWei, eth).Float64()
			cbh.FlashbotsCallBundleResponse.Results[i].CoinbaseDiff = fmt.Sprintf("%.5f Eth", ethCoinbaseDiff)

			bundleGasFeesWei = artemis_eth_units.NewBigFloatFromStr(v.GasFees)
			bundleGasFees, _ = new(big.Float).Quo(bundleGasFeesWei, eth).Float64()
			cbh.FlashbotsCallBundleResponse.Results[i].GasFees = fmt.Sprintf("%.5f Eth", bundleGasFees)
		}

		for i, v := range cbh.Trades {
			if v.AmountInAddr.String() == artemis_trading_constants.WETH9ContractAddress {
				amount := artemis_eth_units.NewBigFloatFromStr(v.AmountIn)
				amountEthF, _ := new(big.Float).Quo(amount, eth).Float64()
				// Format the float to a string with 5 decimal places
				cbh.Trades[i].AmountIn = fmt.Sprintf("%.5f Weth", amountEthF)
			}
			if v.AmountOutAddr.String() == artemis_trading_constants.WETH9ContractAddress {
				amount := artemis_eth_units.NewBigFloatFromStr(v.AmountOut)
				amountEthF, _ := new(big.Float).Quo(amount, eth).Float64()
				// Format the float to a string with 5 decimal places
				cbh.Trades[i].AmountOut = fmt.Sprintf("%.5f Weth", amountEthF)
			}
		}
		cbh.SubmissionTime = ts.ConvertUnixTimeStampToDate(cbh.EventID).String()
		rw = append(rw, cbh)
	}
	return rw, nil
}
