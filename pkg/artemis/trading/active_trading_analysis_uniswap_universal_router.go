package artemis_realtime_trading

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *ActiveTrading) RealTimeProcessUniversalRouterTx(ctx context.Context, tx web3_client.MevTx) ([]web3_client.TradeExecutionFlowJSON, error) {
	subcmd, err := web3_client.NewDecodedUniversalRouterExecCmdFromMap(tx.Args)
	if err != nil {
		return nil, err
	}
	var tfSlice []web3_client.TradeExecutionFlowJSON
	toAddr := tx.Tx.To().String()
	for _, subtx := range subcmd.Commands {
		switch subtx.Command {
		case web3_client.V3SwapExactIn:
			fmt.Println("V3SwapExactIn: ProcessUniversalRouterTxs")
			inputs := subtx.DecodedInputs.(web3_client.V3SwapExactInParams)
			pd, perr := a.u.GetV3PricingData(ctx, inputs.Path)
			if perr != nil {
				a.m.ErrTrackingMetrics.RecordError(web3_client.V3SwapExactIn, pd.V3Pair.PoolAddress)
				log.Err(perr).Msg("V3SwapExactIn: error getting pricing data")
				return nil, perr
			}
			tf := inputs.BinarySearch(pd)
			if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" {
				return nil, errors.New("expectedProfit == 0 or 1")
			}
			tf.Tx = tx.Tx
			tf.InitialPairV3 = pd.V3Pair.ConvertToJSONType()
			tf.Trade.TradeMethod = web3_client.V3SwapExactIn
			a.m.TxFetcherMetrics.TransactionGroup(toAddr, web3_client.V3SwapExactIn)
			a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.Path.TokenIn.String(), inputs.Path.GetEndToken().String())
			a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, web3_client.V3SwapExactIn, pd.V3Pair.PoolAddress, inputs.Path.TokenIn.String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
			tfSlice = append(tfSlice, tf)
		case web3_client.V3SwapExactOut:
			fmt.Println("V3SwapExactOut: ProcessUniversalRouterTxs")
			inputs := subtx.DecodedInputs.(web3_client.V3SwapExactOutParams)
			pd, perr := a.u.GetV3PricingData(ctx, inputs.Path)
			if perr != nil {
				a.m.ErrTrackingMetrics.RecordError(web3_client.V3SwapExactOut, pd.V3Pair.PoolAddress)
				log.Err(perr).Msg("V3SwapExactIn: error getting pricing data")
				return nil, perr
			}
			tf := inputs.BinarySearch(pd)
			tf.Tx = tx.Tx
			if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" {
				return nil, errors.New("expectedProfit == 0 or 1")
			}
			tf.InitialPairV3 = pd.V3Pair.ConvertToJSONType()
			tf.Trade.TradeMethod = web3_client.V3SwapExactOut
			a.m.TxFetcherMetrics.TransactionGroup(toAddr, web3_client.V3SwapExactOut)
			a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.Path.TokenIn.String(), inputs.Path.GetEndToken().String())
			a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, web3_client.V3SwapExactOut, pd.V3Pair.PoolAddress, tf.FrontRunTrade.AmountInAddr.String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
			tfSlice = append(tfSlice, tf)
		case web3_client.V2SwapExactIn:
			fmt.Println("V2SwapExactIn: ProcessUniversalRouterTxs")
			inputs := subtx.DecodedInputs.(web3_client.V2SwapExactInParams)
			pd, perr := a.u.GetV2PricingData(ctx, inputs.Path)
			if perr != nil {
				a.m.ErrTrackingMetrics.RecordError(web3_client.V2SwapExactIn, pd.V2Pair.PairContractAddr)
				log.Err(perr).Msg("V2SwapExactIn: error getting pricing data")
				return nil, perr
			}
			tf := inputs.BinarySearch(pd.V2Pair)
			tf.Tx = tx.Tx
			if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" {
				return nil, errors.New("expectedProfit == 0 or 1")
			}
			tf.Trade.TradeMethod = web3_client.V2SwapExactIn
			tf.InitialPair = pd.V2Pair.ConvertToJSONType()
			a.m.TxFetcherMetrics.TransactionGroup(toAddr, web3_client.V2SwapExactIn)
			pend := len(inputs.Path) - 1
			a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.Path[0].String(), inputs.Path[pend].String())
			a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, web3_client.V2SwapExactIn, pd.V2Pair.PairContractAddr, inputs.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
			tfSlice = append(tfSlice, tf)
		case web3_client.V2SwapExactOut:
			fmt.Println("V2SwapExactOut: ProcessUniversalRouterTxs")
			inputs := subtx.DecodedInputs.(web3_client.V2SwapExactOutParams)
			pd, perr := a.u.GetV2PricingData(ctx, inputs.Path)
			if perr != nil {
				a.m.ErrTrackingMetrics.RecordError(web3_client.V2SwapExactOut, pd.V2Pair.PairContractAddr)
				log.Err(perr).Msg("V2SwapExactOut: error getting pricing data")
				return nil, perr
			}
			tf := inputs.BinarySearch(pd.V2Pair)
			tf.Tx = tx.Tx
			if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" {
				return nil, errors.New("expectedProfit == 0 or 1")
			}
			tf.Trade.TradeMethod = web3_client.V2SwapExactOut
			tf.InitialPair = pd.V2Pair.ConvertToJSONType()
			a.m.TxFetcherMetrics.TransactionGroup(toAddr, web3_client.V2SwapExactOut)
			pend := len(inputs.Path) - 1
			a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.Path[0].String(), inputs.Path[pend].String())
			a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, web3_client.V2SwapExactOut, pd.V2Pair.PairContractAddr, tf.FrontRunTrade.AmountInAddr.String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
			tfSlice = append(tfSlice, tf)
		default:
		}
	}
	return tfSlice, nil
}
