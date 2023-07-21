package web3_client

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/rs/zerolog/log"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
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

func (u *UniswapClient) ProcessUniswapV3RouterTxs(ctx context.Context, tx MevTx, abiFile *abi.ABI) error {
	if strings.HasPrefix(tx.MethodName, multicall) {
		inputs := &Multicall{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			log.Err(err).Msg("failed to decode multicall args")
			return err
		}
		for _, data := range inputs.Data {
			if abiFile == nil {
				mn, args, derr := DecodeTxData(ctx, data, u.MevSmartContractTxMapV3SwapRouterV2.Abi, u.MevSmartContractTxMapV3SwapRouterV2.Filter)
				if derr != nil {
					log.Err(derr).Msg("failed to decode tx data")
					continue
				}
				newTx := tx
				newTx.MethodName = mn
				newTx.Args = args
				err = u.processUniswapV3Txs(ctx, newTx)
				if err != nil {
					log.Err(err).Msg("failed to process uniswap v3 txs")
					continue
				}
			} else {
				mn, args, derr := DecodeTxData(ctx, data, abiFile, nil)
				if derr != nil {
					log.Err(derr).Msg("failed to decode tx data")
					continue
				}
				newTx := tx
				newTx.MethodName = mn
				newTx.Args = args
				err = u.processUniswapV3Txs(ctx, newTx)
				if err != nil {
					log.Err(err).Msg("failed to process uniswap v3 txs")
					continue
				}
			}
		}
	} else {
		err := u.processUniswapV3Txs(ctx, tx)
		if err != nil {
			log.Err(err).Msg("failed to process uniswap v3 txs")
		}
	}
	return nil
}

func (u *UniswapClient) processUniswapV3Txs(ctx context.Context, tx MevTx) error {
	switch tx.MethodName {
	case exactInput:
		inputs := &ExactInputParams{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			log.Err(err).Msg("failed to decode exact input args")
			return err
		}
		pd, perr := u.GetV3PricingData(ctx, inputs.TokenFeePath)
		if perr != nil {
			log.Err(perr).Msg("ExactInput: error getting pricing data")
			return perr
		}
		tf, err := inputs.BinarySearch(pd)
		if err != nil {
			return err
		}
		tfJSON, err := tf.ConvertToJSONType()
		if err != nil {
			log.Err(err).Msg("error converting to json type")
			return err
		}
		fmt.Println("\nsandwich: ==================================ExactInput==================================")
		ts := TradeSummary{
			Tx:            tx,
			Pd:            pd,
			Tf:            tfJSON,
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
			return err
		}
		pd, perr := u.GetV3PricingData(ctx, inputs.TokenFeePath)
		if perr != nil {
			log.Err(perr).Msg("V3SwapExactOut: error getting pricing data")
			return err
		}
		tf, err := inputs.BinarySearch(pd)
		if err != nil {
			return err
		}
		tfJSON, err := tf.ConvertToJSONType()
		if err != nil {
			return err
		}
		fmt.Println("\nsandwich: ==================================ExactOut==================================")
		ts := TradeSummary{
			Tx:            tx,
			Pd:            pd,
			Tf:            tfJSON,
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
			return err
		}
		tfp := artemis_trading_types.TokenFeePath{
			TokenIn: inputs.TokenIn,
			Path: []artemis_trading_types.TokenFee{{
				Token: inputs.TokenOut,
				Fee:   inputs.Fee,
			}},
		}
		inputs.TokenFeePath = tfp
		pd, perr := u.GetV3PricingData(ctx, tfp)
		if perr != nil {
			log.Err(perr).Msg("SwapExactInputSingle: error getting pricing data")
			return err
		}
		tf, err := inputs.BinarySearch(pd)
		if err != nil {
			return err
		}
		tfJSON, err := tf.ConvertToJSONType()
		if err != nil {
			return err
		}
		fmt.Println("\nsandwich: ==================================SwapExactInputSingle==================================")
		ts := TradeSummary{
			Tx:            tx,
			Pd:            pd,
			Tf:            tfJSON,
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
			return err
		}
		tfp := artemis_trading_types.TokenFeePath{
			TokenIn: inputs.TokenIn,
			Path: []artemis_trading_types.TokenFee{{
				Token: inputs.TokenOut,
				Fee:   inputs.Fee,
			}},
		}
		inputs.TokenFeePath = tfp
		pd, perr := u.GetV3PricingData(ctx, tfp)
		if perr != nil {
			log.Err(perr).Msg("SwapExactOutputSingle: error getting pricing data")
			return err
		}
		tf, err := inputs.BinarySearch(pd)
		if err != nil {
			return err
		}
		tfJSON, err := tf.ConvertToJSONType()
		if err != nil {
			return err
		}
		fmt.Println("\nsandwich: ==================================SwapExactOutputSingle==================================")
		ts := TradeSummary{
			Tx:            tx,
			Pd:            pd,
			Tf:            tfJSON,
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
			return err
		}
		pd, err := u.GetV2PricingData(ctx, inputs.Path)
		if err != nil {
			return err
		}
		path := inputs.Path
		tf, err := inputs.BinarySearch(pd.V2Pair)
		if err != nil {
			return err
		}
		tfJSON, err := tf.ConvertToJSONType()
		if err != nil {
			return err
		}
		fmt.Println("\nsandwich: ==================================SwapExactTokensForTokens==================================")
		ts := &TradeSummary{
			Tx:            tx,
			Pd:            pd,
			Tf:            tfJSON,
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
	return nil
}
