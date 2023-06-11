package web3_client

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

// todo add counter

func (u *UniswapClient) ProcessUniversalRouterTxs(ctx context.Context, tx MevTx) {
	subcmd, err := NewDecodedUniversalRouterExecCmdFromMap(tx.Args)
	if err != nil {
		return
	}

	// todo, update this from stub to real
	pair := UniswapV2Pair{}

	// todo needs to save trade analysis results
	for _, subtx := range subcmd.Commands {
		switch subtx.Command {
		case V3SwapExactIn:
			fmt.Println("V3SwapExactIn: ProcessUniversalRouterTxs")
			inputs := subtx.DecodedInputs.(V3SwapExactInParams)
			pd, perr := u.GetPricingData(ctx, inputs.Path.GetPath())
			if perr != nil {
				log.Err(perr).Msg("V3SwapExactIn: error getting pricing data")
				return
			}
			pair = pd.v2Pair
			tf := inputs.BinarySearch(pair)
			tf.InitialPair = pair.ConvertToJSONType()

			fmt.Println("\nsandwich: ==================================V3SwapExactIn==================================")
			ts := TradeSummary{
				Tx:            tx,
				Pd:            pd,
				Tf:            tf,
				TokenAddr:     inputs.Path.TokenIn.String(),
				BuyWithAmount: inputs.AmountIn,
				MinimumAmount: inputs.AmountOutMin,
			}
			u.PrintTradeSummaries(&ts)
			fmt.Println("Sell Token: ", inputs.Path.TokenIn.String(), "Buy Token", inputs.Path.GetEndToken().String(), "Sell BuyWithAmount: ", tf.SandwichPrediction.SellAmount, "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit)
			fmt.Println("sandwich: ====================================V3SwapExactIn==================================")
		case V3SwapExactOut:
			fmt.Println("V3SwapExactOut: ProcessUniversalRouterTxs")
			inputs := subtx.DecodedInputs.(V3SwapExactOutParams)
			pd, perr := u.GetPricingData(ctx, inputs.Path)
			if perr != nil {
				log.Err(perr).Msg("V3SwapExactOut: error getting pricing data")
				return
			}
			pair = pd.v2Pair
			tf := inputs.BinarySearch(pair)
			tf.InitialPair = pair.ConvertToJSONType()
			fmt.Println("\nsandwich: ==================================V3SwapExactOut==================================")
			ts := TradeSummary{
				Tx:            tx,
				Pd:            pd,
				Tf:            tf,
				TokenAddr:     inputs.Path[0].String(),
				BuyWithAmount: inputs.AmountInMax,
				MinimumAmount: inputs.AmountOut,
			}
			u.PrintTradeSummaries(&ts)
			fmt.Println("Sell Token: ", inputs.Path[0].String(), "Buy Token", inputs.Path[1].String(), "Sell BuyWithAmount: ", tf.SandwichPrediction.SellAmount, "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit)
			fmt.Println("sandwich: ====================================V3SwapExactOut==================================")
		case V2SwapExactIn:
			fmt.Println("V2SwapExactIn: ProcessUniversalRouterTxs")
			inputs := subtx.DecodedInputs.(V2SwapExactInParams)
			pd, perr := u.GetPricingData(ctx, inputs.Path)
			if perr != nil {
				log.Err(perr).Msg("V2SwapExactIn: error getting pricing data")
				return
			}
			pair = pd.v2Pair
			tf := inputs.BinarySearch(pair)
			tf.InitialPair = pair.ConvertToJSONType()

			fmt.Println("\nsandwich: ==================================V2SwapExactIn==================================")
			ts := TradeSummary{
				Tx:            tx,
				Pd:            pd,
				Tf:            tf,
				TokenAddr:     inputs.Path[0].String(),
				BuyWithAmount: inputs.AmountIn,
				MinimumAmount: inputs.AmountOutMin,
			}
			u.PrintTradeSummaries(&ts)
			fmt.Println("Sell Token: ", inputs.Path[0].String(), "Buy Token", inputs.Path[1].String(), "Sell BuyWithAmount: ", tf.SandwichPrediction.SellAmount, "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit)
			fmt.Println("sandwich: ====================================V2SwapExactIn==================================")
		case V2SwapExactOut:
			fmt.Println("V2SwapExactOut: ProcessUniversalRouterTxs")

			inputs := subtx.DecodedInputs.(V2SwapExactOutParams)
			pd, perr := u.GetPricingData(ctx, inputs.Path)
			if perr != nil {
				log.Err(perr).Msg("V2SwapExactOut: error getting pricing data")
				return
			}
			pair = pd.v2Pair
			tf := inputs.BinarySearch(pair)
			tf.InitialPair = pair.ConvertToJSONType()

			fmt.Println("\nsandwich: ==================================V2SwapExactOut==================================")
			ts := TradeSummary{
				Tx:            tx,
				Pd:            pd,
				Tf:            tf,
				TokenAddr:     inputs.Path[0].String(),
				BuyWithAmount: inputs.AmountInMax,
				MinimumAmount: inputs.AmountOut,
			}
			u.PrintTradeSummaries(&ts)
			fmt.Println("Sell Token: ", inputs.Path[0].String(), "Buy Token", inputs.Path[1].String(), "Sell BuyWithAmount: ", tf.SandwichPrediction.SellAmount, "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit)
			fmt.Println("sandwich: ====================================V2SwapExactOut==================================")
		default:
		}
	}
}
