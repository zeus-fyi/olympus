package artemis_realtime_trading

import (
	"context"
	"fmt"

	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *ActiveTrading) RealTimeProcessUniversalRouterTx(ctx context.Context, tx web3_client.MevTx) {
	subcmd, err := web3_client.NewDecodedUniversalRouterExecCmdFromMap(tx.Args)
	if err != nil {
		return
	}
	if tx.Tx.To() == nil {
		return
	}
	toAddr := tx.Tx.To().String()
	for _, subtx := range subcmd.Commands {
		switch subtx.Command {
		case web3_client.V3SwapExactIn:
			fmt.Println("V3SwapExactIn: ProcessUniversalRouterTxs")
			inputs := subtx.DecodedInputs.(web3_client.V3SwapExactInParams)
			a.m.TxFetcherMetrics.TransactionGroup(toAddr, web3_client.V3SwapExactIn)
			a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.Path.TokenIn.String(), inputs.Path.GetEndToken().String())
			//pd, perr := u.GetV3PricingData(ctx, inputs.Path)
			//if perr != nil {
			//	log.Err(perr).Msg("V3SwapExactIn: error getting pricing data")
			//	return
			//}
			//tf := inputs.BinarySearch(pd)
			//tf.Trade.TradeMethod = web3_client.V3SwapExactIn
		case web3_client.V3SwapExactOut:
			fmt.Println("V3SwapExactOut: ProcessUniversalRouterTxs")
			inputs := subtx.DecodedInputs.(web3_client.V3SwapExactOutParams)
			a.m.TxFetcherMetrics.TransactionGroup(toAddr, web3_client.V3SwapExactOut)
			a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.Path.TokenIn.String(), inputs.Path.GetEndToken().String())
			//pd, perr := u.GetV3PricingData(ctx, inputs.Path)
			//if perr != nil {
			//	log.Err(perr).Msg("V3SwapExactOut: error getting pricing data")
			//	return
			//}
			//tf := inputs.BinarySearch(pd)
			//tf.InitialPairV3 = pd.V3Pair.ConvertToJSONType()
			//ts := web3_client.TradeSummary{
			//	Tx:            tx,
			//	Pd:            pd,
			//	Tf:            tf,
			//	TokenAddr:     inputs.Path.TokenIn.String(),
			//	BuyWithAmount: inputs.AmountInMax,
			//	MinimumAmount: inputs.AmountOut,
			//}
			//tf.Trade.TradeMethod = web3_client.V3SwapExactOut
		case web3_client.V2SwapExactIn:
			fmt.Println("V2SwapExactIn: ProcessUniversalRouterTxs")
			inputs := subtx.DecodedInputs.(web3_client.V2SwapExactInParams)
			a.m.TxFetcherMetrics.TransactionGroup(toAddr, web3_client.V2SwapExactIn)
			pend := len(inputs.Path) - 1
			a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.Path[0].String(), inputs.Path[pend].String())
			//pd, perr := u.GetV2PricingData(ctx, inputs.Path)
			//if perr != nil {
			//	log.Err(perr).Msg("V2SwapExactIn: error getting pricing data")
			//	return
			//}
			//pair := pd.V2Pair
			//tf := inputs.BinarySearch(pair)
			//ts := web3_client.TradeSummary{
			//	Tx:            tx,
			//	Pd:            pd,
			//	Tf:            tf,
			//	TokenAddr:     inputs.Path[0].String(),
			//	BuyWithAmount: inputs.AmountIn,
			//	MinimumAmount: inputs.AmountOutMin,
			//}
			//tf.Trade.TradeMethod =  web3_client.V2SwapExactIn
		case web3_client.V2SwapExactOut:
			fmt.Println("V2SwapExactOut: ProcessUniversalRouterTxs")
			inputs := subtx.DecodedInputs.(web3_client.V2SwapExactOutParams)
			a.m.TxFetcherMetrics.TransactionGroup(toAddr, web3_client.V2SwapExactOut)
			pend := len(inputs.Path) - 1
			a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.Path[0].String(), inputs.Path[pend].String())
			//pd, perr := u.GetV2PricingData(ctx, inputs.Path)
			//if perr != nil {
			//	log.Err(perr).Msg("V2SwapExactOut: error getting pricing data")
			//	return
			//}
			//pair := pd.V2Pair
			//tf := inputs.BinarySearch(pair)
			//tf.InitialPair = pair.ConvertToJSONType()
			//ts := web3_client.TradeSummary{
			//	Tx:            tx,
			//	Pd:            pd,
			//	Tf:            tf,
			//	TokenAddr:     inputs.Path[0].String(),
			//	BuyWithAmount: inputs.AmountInMax,
			//	MinimumAmount: inputs.AmountOut,
			//}
			//a.m.TxFetcherMetrics.TransactionCurrencyIn(tx.Tx.To().String(), inputs.Path[0].String(), inputs.AmountInMax)
			//tf.Trade.TradeMethod =  web3_client.V2SwapExactOut
		default:
		}
	}
}
