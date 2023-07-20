package web3_client

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/rs/zerolog/log"
)

// todo fix this
// {"level":"error","error":"pair address length is not 2","time":1687224704,"message":"V2SwapExactIn: error getting pricing data"}

func (u *UniswapClient) ProcessUniversalRouterTxs(ctx context.Context, tx MevTx, abiFile *abi.ABI) error {
	subcmd, serr := NewDecodedUniversalRouterExecCmdFromMap(tx.Args, abiFile)
	if serr != nil {
		return serr
	}

	// todo needs to compound all trades per execution command
	count := 0
	for _, subtx := range subcmd.Commands {
		switch subtx.Command {
		case V3SwapExactIn:
			fmt.Println("V3SwapExactIn: ProcessUniversalRouterTxs")
			inputs := subtx.DecodedInputs.(V3SwapExactInParams)
			pd, perr := u.GetV3PricingData(ctx, inputs.Path)
			if perr != nil {
				log.Err(perr).Msg("V3SwapExactIn: error getting pricing data")
				return perr
			}
			tf, err := inputs.BinarySearch(pd)
			if err != nil {
				log.Err(err).Msg("V3SwapExactIn: error getting binary search")
				return err
			}
			tf.InitialPairV3 = pd.V3Pair.ConvertToJSONType()
			fmt.Println("\nsandwich: ==================================V3SwapExactIn==================================")
			ts := TradeSummary{
				Tx:            tx,
				Pd:            pd,
				Tf:            tf,
				TokenAddr:     inputs.Path.TokenIn.String(),
				BuyWithAmount: inputs.AmountIn,
				MinimumAmount: inputs.AmountOutMin,
			}
			tf.Trade.TradeMethod = V3SwapExactIn
			u.PrintTradeSummaries(&ts)
			fmt.Println("txHash: ", tx.Tx.Hash().String())
			fmt.Println("Sell Token: ", inputs.Path.TokenIn.String(), "Buy Token", inputs.Path.GetEndToken().String(), "Sell Amount: ", tf.SandwichPrediction.SellAmount, "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit)
			fmt.Println("sandwich: ====================================V3SwapExactIn==================================")
			count++
		case V3SwapExactOut:
			fmt.Println("V3SwapExactOut: ProcessUniversalRouterTxs")
			inputs := subtx.DecodedInputs.(V3SwapExactOutParams)
			pd, perr := u.GetV3PricingData(ctx, inputs.Path)
			if perr != nil {
				log.Err(perr).Msg("V3SwapExactOut: error getting pricing data")
				return perr
			}
			tf, err := inputs.BinarySearch(pd)
			if err != nil {
				log.Err(err).Msg("V3SwapExactOut: error getting binary search")
				return err
			}
			tf.InitialPairV3 = pd.V3Pair.ConvertToJSONType()
			fmt.Println("\nsandwich: ==================================V3SwapExactOut==================================")
			ts := TradeSummary{
				Tx:            tx,
				Pd:            pd,
				Tf:            tf,
				TokenAddr:     inputs.Path.TokenIn.String(),
				BuyWithAmount: inputs.AmountInMax,
				MinimumAmount: inputs.AmountOut,
			}
			tf.Trade.TradeMethod = V3SwapExactOut
			u.PrintTradeSummaries(&ts)
			fmt.Println("txHash: ", tx.Tx.Hash().String())
			fmt.Println("Sell Token: ", inputs.Path.TokenIn.String(), "Buy Token", inputs.Path.GetEndToken().String(), "Sell Amount: ", tf.SandwichPrediction.SellAmount, "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit)
			fmt.Println("sandwich: ====================================V3SwapExactOut==================================")
			count++
		case V2SwapExactIn:
			fmt.Println("V2SwapExactIn: ProcessUniversalRouterTxs")
			inputs := subtx.DecodedInputs.(V2SwapExactInParams)
			pd, perr := u.GetV2PricingData(ctx, inputs.Path)
			if perr != nil {
				log.Err(perr).Msg("V2SwapExactIn: error getting pricing data")
				return perr
			}
			pair := pd.V2Pair
			tf, err := inputs.BinarySearch(pair)
			if err != nil {
				log.Err(err).Msg("V2SwapExactIn: error getting binary search")
				return err
			}
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
			tf.Trade.TradeMethod = V2SwapExactIn
			u.PrintTradeSummaries(&ts)
			fmt.Println("txHash: ", tx.Tx.Hash().String())
			fmt.Println("Sell Token: ", inputs.Path[0].String(), "Buy Token", inputs.Path[1].String(), "Sell Amount: ", tf.SandwichPrediction.SellAmount, "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit)
			fmt.Println("sandwich: ====================================V2SwapExactIn==================================")
			count++
		case V2SwapExactOut:
			fmt.Println("V2SwapExactOut: ProcessUniversalRouterTxs")
			inputs := subtx.DecodedInputs.(V2SwapExactOutParams)
			pd, perr := u.GetV2PricingData(ctx, inputs.Path)
			if perr != nil {
				log.Err(perr).Msg("V2SwapExactOut: error getting pricing data")
				return perr
			}
			pair := pd.V2Pair
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
			tf.Trade.TradeMethod = V2SwapExactOut
			u.PrintTradeSummaries(&ts)
			fmt.Println("txHash: ", tx.Tx.Hash().String())
			fmt.Println("Sell Token: ", inputs.Path[0].String(), "Buy Token", inputs.Path[1].String(), "Sell Amount: ", tf.SandwichPrediction.SellAmount, "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit)
			fmt.Println("sandwich: ====================================V2SwapExactOut==================================")
			count++
		default:
		}
	}
	fmt.Println("filtered total trades: ", count)
	return nil
}
