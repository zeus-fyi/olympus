package web3_client

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	swapExactInputSingle    = "swapExactInputSingle"
	swapExactOutputSingle   = "swapExactOutputSingle"
	swapExactInputMultihop  = "swapExactInputMultihop"
	swapExactOutputMultihop = "swapExactOutputMultihop"
	exactInput              = "exactInput"
	exactOutput             = "exactOutput"
	multicall               = "multicall"
	multicall0              = "multicall0"
	multicall1              = "multicall1"
)

/*
UniswapV2 — Router: can perform basic ERC-20 swaps
UniswapV2 — Router2: can perform basic ERC-20 plus fee on transfer swaps
UniswapV3 — Router: can perform ERC-20 swaps of any kind, limited to UniswapV3 pools
UniswapV3 — Router2: can perform ERC-20 swaps of any kind through both UniswapV2 and UniswapV3 pools
*/

func (u *UniswapClient) ProcessUniswapV3RouterTxs(ctx context.Context, tx MevTx) {
	if strings.HasPrefix(tx.MethodName, multicall) {
		inputs := &Multicall{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			log.Err(err).Msg("failed to decode multicall args")
			return
		}
		for _, data := range inputs.Data {
			mn, args, derr := DecodeTxData(ctx, data, u.MevSmartContractTxMapV3)
			if derr != nil {
				log.Err(derr).Msg("failed to decode tx data")
				continue
			}
			newTx := tx
			newTx.MethodName = mn
			newTx.Args = args
			u.processUniswapV3Txs(ctx, newTx)
		}
	} else {
		u.processUniswapV3Txs(ctx, tx)
	}
	return
}

func (u *UniswapClient) processUniswapV3Txs(ctx context.Context, tx MevTx) {
	switch tx.MethodName {
	case exactInput:
		inputs := &ExactInputParams{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			log.Err(err).Msg("failed to decode exact input args")
			return
		}
		pd, perr := u.GetV3PricingData(ctx, inputs.TokenFeePath)
		if perr != nil {
			log.Err(perr).Msg("ExactInput: error getting pricing data")
			return
		}
		tf := inputs.BinarySearch(pd)
		tf.Trade.TradeMethod = exactInput
		tf.InitialPairV3 = pd.v3Pair.ConvertToJSONType()
		fmt.Println("\nsandwich: ==================================ExactInput==================================")
		ts := TradeSummary{
			Tx:            tx,
			Pd:            pd,
			Tf:            tf,
			TokenAddr:     inputs.TokenFeePath.TokenIn.String(),
			BuyWithAmount: inputs.AmountIn,
			MinimumAmount: inputs.AmountOutMinimum,
		}
		u.PrintTradeSummaries(&ts)
		fmt.Println("txHash: ", tx.Tx.Hash().String())
		fmt.Println("Sell Token: ", inputs.TokenFeePath.TokenIn.String(), "Buy Token", inputs.TokenFeePath.GetEndToken().String(), "Sell Amount: ", tf.SandwichPrediction.SellAmount, "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit)
		fmt.Println("sandwich: ====================================ExactInput==================================")
	case exactOutput:
		inputs := &ExactOutputParams{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			log.Err(err).Msg("failed to decode exact output args")
			return
		}
		pd, perr := u.GetV3PricingData(ctx, inputs.TokenFeePath)
		if perr != nil {
			log.Err(perr).Msg("V3SwapExactOut: error getting pricing data")
			return
		}
		tf := inputs.BinarySearch(pd)
		tf.Trade.TradeMethod = exactOutput
		tf.InitialPairV3 = pd.v3Pair.ConvertToJSONType()
		fmt.Println("\nsandwich: ==================================ExactOut==================================")
		ts := TradeSummary{
			Tx:            tx,
			Pd:            pd,
			Tf:            tf,
			TokenAddr:     inputs.TokenFeePath.TokenIn.String(),
			BuyWithAmount: inputs.AmountInMaximum,
			MinimumAmount: inputs.AmountOut,
		}
		u.PrintTradeSummaries(&ts)
		fmt.Println("txHash: ", tx.Tx.Hash().String())
		fmt.Println("Sell Token: ", inputs.TokenFeePath.TokenIn.String(), "Buy Token", inputs.TokenFeePath.GetEndToken().String(), "Sell Amount: ", tf.SandwichPrediction.SellAmount, "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit)
		fmt.Println("sandwich: ====================================ExactOut==================================")
	case swapExactInputSingle:
		inputs := &SwapExactInputSingleArgs{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			log.Err(err).Msg("failed to decode swap exact input single args")
			return
		}
		tfp := TokenFeePath{
			TokenIn: inputs.TokenIn,
			Path: []TokenFee{{
				Token: inputs.TokenOut,
				Fee:   inputs.Fee,
			}},
		}
		inputs.TokenFeePath = tfp
		pd, perr := u.GetV3PricingData(ctx, tfp)
		if perr != nil {
			log.Err(perr).Msg("SwapExactInputSingle: error getting pricing data")
			return
		}
		tf := inputs.BinarySearch(pd)
		tf.Trade.TradeMethod = swapExactInputSingle
		tf.InitialPairV3 = pd.v3Pair.ConvertToJSONType()
		fmt.Println("\nsandwich: ==================================SwapExactInputSingle==================================")
		ts := TradeSummary{
			Tx:            tx,
			Pd:            pd,
			Tf:            tf,
			TokenAddr:     inputs.TokenFeePath.TokenIn.String(),
			BuyWithAmount: inputs.AmountIn,
			MinimumAmount: inputs.AmountOutMinimum,
		}
		u.PrintTradeSummaries(&ts)
		fmt.Println("txHash: ", tx.Tx.Hash().String())
		fmt.Println("Sell Token: ", inputs.TokenFeePath.TokenIn.String(), "Buy Token", inputs.TokenFeePath.GetEndToken().String(), "Sell Amount: ", tf.SandwichPrediction.SellAmount, "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit)
		fmt.Println("sandwich: ====================================SwapExactInputSingle==================================")
	case swapExactOutputSingle:
		inputs := &SwapExactOutputSingleArgs{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			log.Err(err).Msg("failed to decode swap exact output single args")
			return
		}
		tfp := TokenFeePath{
			TokenIn: inputs.TokenIn,
			Path: []TokenFee{{
				Token: inputs.TokenOut,
				Fee:   inputs.Fee,
			}},
		}
		inputs.TokenFeePath = tfp
		pd, perr := u.GetV3PricingData(ctx, tfp)
		if perr != nil {
			log.Err(perr).Msg("SwapExactOutputSingle: error getting pricing data")
			return
		}
		tf := inputs.BinarySearch(pd)
		tf.Trade.TradeMethod = swapExactOutputSingle
		tf.InitialPairV3 = pd.v3Pair.ConvertToJSONType()
		fmt.Println("\nsandwich: ==================================SwapExactOutputSingle==================================")
		ts := TradeSummary{
			Tx:            tx,
			Pd:            pd,
			Tf:            tf,
			TokenAddr:     inputs.TokenFeePath.TokenIn.String(),
			BuyWithAmount: inputs.AmountInMaximum,
			MinimumAmount: inputs.AmountOut,
		}
		u.PrintTradeSummaries(&ts)
		fmt.Println("txHash: ", tx.Tx.Hash().String())
		fmt.Println("Sell Token: ", inputs.TokenFeePath.TokenIn.String(), "Buy Token", inputs.TokenFeePath.GetEndToken().String(), "Sell Amount: ", tf.SandwichPrediction.SellAmount, "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit)
		fmt.Println("sandwich: ====================================SwapExactOutputSingle==================================")
	case swapExactTokensForTokens:
		inputs := &SwapExactTokensForTokensParamsV3{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			log.Err(err).Msg("swapExactTokensForTokens: failed to decode swap exact tokens for tokens args")
			return
		}
		pd, err := u.GetPricingData(ctx, inputs.Path)
		if err != nil {
			return
		}
		path := inputs.Path
		initialPair := pd.v2Pair
		tf := inputs.BinarySearch(pd.v2Pair)
		tf.InitialPair = initialPair.ConvertToJSONType()
		fmt.Println("\nsandwich: ==================================SwapExactTokensForTokens==================================")
		ts := &TradeSummary{
			Tx:            tx,
			Pd:            pd,
			Tf:            tf,
			TokenAddr:     path[0].String(),
			BuyWithAmount: inputs.AmountIn,
			MinimumAmount: inputs.AmountOutMin,
		}
		u.PrintTradeSummaries(ts)
		fmt.Println("txHash: ", tx.Tx.Hash().String())
		fmt.Println("Sell Token: ", path[0].String(), "Buy Token", path[1].String(), "SandwichPrediction Sell Amount: ", tf.SandwichPrediction.SellAmount, "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit)
		fmt.Println("sandwich: ====================================SwapExactTokensForTokens==================================")
	case swapExactInputMultihop:
	case swapExactOutputMultihop:
	}
}
