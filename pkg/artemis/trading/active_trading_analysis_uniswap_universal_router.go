package artemis_realtime_trading

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *ActiveTrading) RealTimeProcessUniversalRouterTx(ctx context.Context, tx web3_client.MevTx) ([]web3_client.TradeExecutionFlowJSON, error) {
	if tx.Tx == nil || tx.Args == nil {
		return nil, errors.New("tx is nil")
	}
	subcmd, err := web3_client.NewDecodedUniversalRouterExecCmdFromMap(tx.Args)
	if err != nil {
		return nil, err
	}
	bn, err := artemis_trading_cache.GetLatestBlock(ctx)
	if err != nil {
		log.Err(err).Msg("failed to get latest block")
		return nil, errors.New("ailed to get latest block")
	}
	var tfSlice []web3_client.TradeExecutionFlowJSON
	toAddr := tx.Tx.To().String()
	for _, subtx := range subcmd.Commands {
		switch subtx.Command {
		case web3_client.V3SwapExactIn:
			//fmt.Println("V3SwapExactIn: ProcessUniversalRouterTxs")
			inputs := subtx.DecodedInputs.(web3_client.V3SwapExactInParams)
			pd, perr := a.GetUniswapClient().GetV3PricingData(ctx, inputs.Path)
			if perr != nil {
				if pd != nil {
					a.GetMetricsClient().ErrTrackingMetrics.RecordError(web3_client.V3SwapExactIn, pd.V3Pair.PoolAddress)
				}
				//log.Err(perr).Msg("V3SwapExactIn: error getting pricing data")
				return nil, perr
			}
			tf := inputs.BinarySearch(pd)
			if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" || tf.SandwichPrediction.ExpectedProfit == "" {
				return nil, errors.New("expectedProfit == 0 or 1")
			}
			newJsonTx := artemis_trading_types.JSONTx{}
			err = newJsonTx.UnmarshalTx(tx.Tx)
			if err != nil {
				return nil, err
			}
			tf.Tx = newJsonTx
			tf.InitialPairV3 = pd.V3Pair.ConvertToJSONType()
			tf.Trade.TradeMethod = web3_client.V3SwapExactIn
			a.GetMetricsClient().TxFetcherMetrics.TransactionGroup(toAddr, web3_client.V3SwapExactIn)
			a.GetMetricsClient().TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.Path.TokenIn.String(), inputs.Path.GetEndToken().String())
			a.GetMetricsClient().TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, web3_client.V3SwapExactIn, pd.V3Pair.PoolAddress, inputs.Path.TokenIn.String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
			tfSlice = append(tfSlice, tf)

			log.Info().Msg("saving mempool tx")
			err = a.SaveMempoolTx(ctx, bn, []web3_client.TradeExecutionFlowJSON{tf})
			if err != nil {
				log.Err(err).Msg("failed to save mempool tx")
				return nil, errors.New("failed to save mempool tx")
			}
		case web3_client.V3SwapExactOut:
			//fmt.Println("V3SwapExactOut: ProcessUniversalRouterTxs")
			inputs := subtx.DecodedInputs.(web3_client.V3SwapExactOutParams)
			pd, perr := a.GetUniswapClient().GetV3PricingData(ctx, inputs.Path)
			if perr != nil {
				if pd != nil {
					a.GetMetricsClient().ErrTrackingMetrics.RecordError(web3_client.V3SwapExactOut, pd.V3Pair.PoolAddress)
				}
				//log.Err(perr).Msg("V3SwapExactIn: error getting pricing data")
				return nil, perr
			}
			tf := inputs.BinarySearch(pd)
			newJsonTx := artemis_trading_types.JSONTx{}
			err = newJsonTx.UnmarshalTx(tx.Tx)
			if err != nil {
				return nil, err
			}
			tf.Tx = newJsonTx
			if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" || tf.SandwichPrediction.ExpectedProfit == "" {
				return nil, errors.New("expectedProfit == 0 or 1")
			}
			tf.InitialPairV3 = pd.V3Pair.ConvertToJSONType()
			tf.Trade.TradeMethod = web3_client.V3SwapExactOut
			a.GetMetricsClient().TxFetcherMetrics.TransactionGroup(toAddr, web3_client.V3SwapExactOut)
			a.GetMetricsClient().TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.Path.TokenIn.String(), inputs.Path.GetEndToken().String())
			a.GetMetricsClient().TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, web3_client.V3SwapExactOut, pd.V3Pair.PoolAddress, tf.FrontRunTrade.AmountInAddr.String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
			tfSlice = append(tfSlice, tf)

			log.Info().Msg("saving mempool tx")
			err = a.SaveMempoolTx(ctx, bn, []web3_client.TradeExecutionFlowJSON{tf})
			if err != nil {
				log.Err(err).Msg("failed to save mempool tx")
				return nil, errors.New("failed to save mempool tx")
			}
		case web3_client.V2SwapExactIn:
			//fmt.Println("V2SwapExactIn: ProcessUniversalRouterTxs")
			inputs := subtx.DecodedInputs.(web3_client.V2SwapExactInParams)
			pd, perr := a.GetUniswapClient().GetV2PricingData(ctx, inputs.Path)
			if perr != nil {
				if pd != nil {
					a.GetMetricsClient().ErrTrackingMetrics.RecordError(web3_client.V2SwapExactIn, pd.V2Pair.PairContractAddr)
				}
				//log.Err(perr).Msg("V2SwapExactIn: error getting pricing data")
				return nil, perr
			}
			tf := inputs.BinarySearch(pd.V2Pair)
			newJsonTx := artemis_trading_types.JSONTx{}
			err = newJsonTx.UnmarshalTx(tx.Tx)
			if err != nil {
				return nil, err
			}
			tf.Tx = newJsonTx
			if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" || tf.SandwichPrediction.ExpectedProfit == "" {
				return nil, errors.New("expectedProfit == 0 or 1")
			}
			tf.Trade.TradeMethod = web3_client.V2SwapExactIn
			tf.InitialPair = pd.V2Pair.ConvertToJSONType()
			a.GetMetricsClient().TxFetcherMetrics.TransactionGroup(toAddr, web3_client.V2SwapExactIn)
			pend := len(inputs.Path) - 1
			a.GetMetricsClient().TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.Path[0].String(), inputs.Path[pend].String())
			a.GetMetricsClient().TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, web3_client.V2SwapExactIn, pd.V2Pair.PairContractAddr, inputs.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
			tfSlice = append(tfSlice, tf)
			log.Info().Msg("saving mempool tx")
			err = a.SaveMempoolTx(ctx, bn, []web3_client.TradeExecutionFlowJSON{tf})
			if err != nil {
				log.Err(err).Msg("failed to save mempool tx")
				return nil, errors.New("failed to save mempool tx")
			}
		case web3_client.V2SwapExactOut:
			//fmt.Println("V2SwapExactOut: ProcessUniversalRouterTxs")
			inputs := subtx.DecodedInputs.(web3_client.V2SwapExactOutParams)
			pd, perr := a.GetUniswapClient().GetV2PricingData(ctx, inputs.Path)
			if perr != nil {
				if pd != nil {
					a.GetMetricsClient().ErrTrackingMetrics.RecordError(web3_client.V2SwapExactOut, pd.V2Pair.PairContractAddr)
				}
				//log.Err(perr).Msg("V2SwapExactOut: error getting pricing data")
				return nil, perr
			}
			tf := inputs.BinarySearch(pd.V2Pair)
			newJsonTx := artemis_trading_types.JSONTx{}
			err = newJsonTx.UnmarshalTx(tx.Tx)
			if err != nil {
				return nil, err
			}
			tf.Tx = newJsonTx
			if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" || tf.SandwichPrediction.ExpectedProfit == "" {
				return nil, errors.New("expectedProfit == 0 or 1")
			}
			tf.Trade.TradeMethod = web3_client.V2SwapExactOut
			tf.InitialPair = pd.V2Pair.ConvertToJSONType()
			a.GetMetricsClient().TxFetcherMetrics.TransactionGroup(toAddr, web3_client.V2SwapExactOut)
			pend := len(inputs.Path) - 1
			a.GetMetricsClient().TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.Path[0].String(), inputs.Path[pend].String())
			a.GetMetricsClient().TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, web3_client.V2SwapExactOut, pd.V2Pair.PairContractAddr, tf.FrontRunTrade.AmountInAddr.String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
			tfSlice = append(tfSlice, tf)
			log.Info().Msg("saving mempool tx")
			err = a.SaveMempoolTx(ctx, bn, []web3_client.TradeExecutionFlowJSON{tf})
			if err != nil {
				log.Err(err).Msg("failed to save mempool tx")
				return nil, errors.New("failed to save mempool tx")
			}
		default:
		}
	}
	return tfSlice, nil
}
