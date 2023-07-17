package artemis_realtime_trading

import (
	"context"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/rs/zerolog/log"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
)

const (
	multicall               = "multicall"
	swapExactInputSingle    = "swapExactInputSingle"
	swapExactOutputSingle   = "swapExactOutputSingle"
	swapExactInputMultihop  = "swapExactInputMultihop"
	swapExactOutputMultihop = "swapExactOutputMultihop"
	exactInput              = "exactInput"
	exactOutput             = "exactOutput"
)

func (a *ActiveTrading) RealTimeProcessUniswapV3RouterTx(ctx context.Context, tx web3_client.MevTx, abiFile *abi.ABI, filter *strings_filter.FilterOpts) ([]web3_client.TradeExecutionFlowJSON, error) {
	w3a := a.GetUniswapClient().Web3Client.Web3Actions
	toAddr := tx.Tx.To().String()
	var tfSlice []web3_client.TradeExecutionFlowJSON
	if strings.HasPrefix(tx.MethodName, multicall) {
		a.GetMetricsClient().TxFetcherMetrics.TransactionGroup(toAddr, multicall)
		inputs := &web3_client.Multicall{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			log.Err(err).Msg("failed to decode multicall args")
			return nil, err
		}
		for _, data := range inputs.Data {
			mn, args, derr := web3_client.DecodeTxData(ctx, data, abiFile, filter)
			if derr != nil {
				log.Err(derr).Msg("failed to decode tx data")
				continue
			}
			newTx := tx
			newTx.MethodName = mn
			newTx.Args = args
			tf, terr := a.processUniswapV3Txs(ctx, newTx)
			if terr != nil {
				log.Err(terr).Msg("failed to process uniswap v3 tx")
				continue
			}
			tfSlice = append(tfSlice, tf...)
		}
	} else {
		tf, err := a.processUniswapV3Txs(ctx, tx)
		if err != nil {
			log.Err(err).Msg("failed to process uniswap v3 tx")
			return nil, err
		}
		tfSlice = append(tfSlice, tf...)
	}
	return tfSlice, nil
}

func (a *ActiveTrading) processUniswapV3Txs(ctx context.Context, tx web3_client.MevTx) ([]web3_client.TradeExecutionFlowJSON, error) {
	bn, berr := artemis_trading_cache.GetLatestBlock(ctx)
	if berr != nil {
		log.Err(berr).Msg("failed to get latest block")
		return nil, errors.New("ailed to get latest block")
	}
	var tfSlice []web3_client.TradeExecutionFlowJSON
	toAddr := tx.Tx.To().String()
	switch tx.MethodName {
	case exactInput:
		inputs := &web3_client.ExactInputParams{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			log.Err(err).Msg("failed to decode exact input args")
			return nil, err
		}
		pd, err := a.GetUniswapClient().GetV3PricingData(ctx, inputs.TokenFeePath)
		if err != nil {
			if pd != nil {
				a.GetMetricsClient().ErrTrackingMetrics.RecordError(exactInput, pd.V3Pair.PoolAddress)
			}
			//log.Err(err).Msg("failed to get pricing data")
			return nil, err
		}
		if pd == nil {
			return nil, errors.New("pd is nil")
		}
		tf := inputs.BinarySearch(pd)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" || tf.SandwichPrediction.ExpectedProfit == "" {
			return nil, errors.New("expectedProfit == 0 or 1")
		}
		tf.Trade.TradeMethod = exactInput
		newTx := artemis_trading_types.JSONTx{}
		err = newTx.UnmarshalTx(tx.Tx)
		if err != nil {
			log.Err(err).Msg("failed to unmarshal tx")
			return nil, err
		}
		tf.Tx = newTx
		tf.InitialPairV3 = pd.V3Pair.ConvertToJSONType()
		a.GetMetricsClient().TxFetcherMetrics.TransactionGroup(toAddr, exactInput)
		a.GetMetricsClient().TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.TokenFeePath.TokenIn.String(), inputs.TokenFeePath.GetEndToken().String())
		a.GetMetricsClient().TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, exactInput, pd.V3Pair.PoolAddress, inputs.TokenFeePath.TokenIn.String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		a.GetMetricsClient().StageProgressionMetrics.CountPostProcessTx(float64(1))
		tfSlice = append(tfSlice, tf)
		log.Info().Msg("saving mempool tx")
		err = a.SaveMempoolTx(ctx, bn, []web3_client.TradeExecutionFlowJSON{tf})
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
	case exactOutput:
		inputs := &web3_client.ExactOutputParams{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			log.Err(err).Msg("failed to decode exact output args")
			return nil, err
		}
		pd, err := a.GetUniswapClient().GetV3PricingData(ctx, inputs.TokenFeePath)
		if err != nil {
			if pd != nil {
				a.GetMetricsClient().ErrTrackingMetrics.RecordError(exactOutput, pd.V3Pair.PoolAddress)
			}
			//log.Err(err).Msg("failed to get pricing data")
			return nil, err
		}
		if pd == nil {
			return nil, errors.New("pd is nil")
		}
		tf := inputs.BinarySearch(pd)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" || tf.SandwichPrediction.ExpectedProfit == "" {
			return nil, errors.New("expectedProfit == 0 or 1")
		}
		tf.Trade.TradeMethod = exactOutput
		newTx := artemis_trading_types.JSONTx{}
		err = newTx.UnmarshalTx(tx.Tx)
		if err != nil {
			log.Err(err).Msg("failed to unmarshal tx")
			return nil, err
		}
		tf.Tx = newTx
		tf.InitialPairV3 = pd.V3Pair.ConvertToJSONType()
		a.GetMetricsClient().TxFetcherMetrics.TransactionGroup(toAddr, exactOutput)
		a.GetMetricsClient().TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.TokenFeePath.TokenIn.String(), inputs.TokenFeePath.GetEndToken().String())
		a.GetMetricsClient().TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, exactOutput, pd.V3Pair.PoolAddress, tf.FrontRunTrade.AmountInAddr.String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		a.GetMetricsClient().StageProgressionMetrics.CountPostProcessTx(float64(1))

		tfSlice = append(tfSlice, tf)
		log.Info().Msg("saving mempool tx")
		err = a.SaveMempoolTx(ctx, bn, []web3_client.TradeExecutionFlowJSON{tf})
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
	case swapExactInputSingle:
		inputs := &web3_client.SwapExactInputSingleArgs{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			log.Err(err).Msg("failed to decode swap exact input single args")
			return nil, err
		}
		pd, err := a.GetUniswapClient().GetV3PricingData(ctx, inputs.TokenFeePath)
		if err != nil {
			if pd != nil {
				a.GetMetricsClient().ErrTrackingMetrics.RecordError(swapExactInputSingle, pd.V3Pair.PoolAddress)
			}
			//log.Err(err).Msg("failed to get pricing data")
			return nil, err
		}
		if pd == nil {
			return nil, errors.New("pd is nil")
		}
		tf := inputs.BinarySearch(pd)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" || tf.SandwichPrediction.ExpectedProfit == "" {
			return nil, errors.New("expectedProfit == 0 or 1")
		}
		tf.Trade.TradeMethod = swapExactInputSingle
		newTx := artemis_trading_types.JSONTx{}
		err = newTx.UnmarshalTx(tx.Tx)
		if err != nil {
			//log.Err(err).Msg("failed to unmarshal tx")
			return nil, err
		}
		tf.Tx = newTx
		tf.InitialPairV3 = pd.V3Pair.ConvertToJSONType()
		a.GetMetricsClient().TxFetcherMetrics.TransactionGroup(toAddr, swapExactInputSingle)
		a.GetMetricsClient().TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.TokenFeePath.TokenIn.String(), inputs.TokenFeePath.GetEndToken().String())
		a.GetMetricsClient().TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactInputSingle, pd.V3Pair.PoolAddress, inputs.TokenFeePath.TokenIn.String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		a.GetMetricsClient().StageProgressionMetrics.CountPostProcessTx(float64(1))
		tfSlice = append(tfSlice, tf)
		log.Info().Msg("saving mempool tx")
		err = a.SaveMempoolTx(ctx, bn, []web3_client.TradeExecutionFlowJSON{tf})
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
	case swapExactOutputSingle:
		inputs := &web3_client.SwapExactOutputSingleArgs{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			//log.Err(err).Msg("failed to decode swap exact output single args")
			return nil, err
		}
		pd, err := a.GetUniswapClient().GetV3PricingData(ctx, inputs.TokenFeePath)
		if err != nil {
			if pd != nil {
				a.GetMetricsClient().ErrTrackingMetrics.RecordError(swapExactOutputSingle, pd.V3Pair.PoolAddress)
			}
			//log.Err(err).Msg("failed to get pricing data")
			return nil, err
		}
		if pd == nil {
			return nil, errors.New("pd is nil")
		}
		tf := inputs.BinarySearch(pd)
		newTx := artemis_trading_types.JSONTx{}
		err = newTx.UnmarshalTx(tx.Tx)
		if err != nil {
			//log.Err(err).Msg("failed to unmarshal tx")
			return nil, err
		}
		tf.Tx = newTx
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" || tf.SandwichPrediction.ExpectedProfit == "" {
			return nil, errors.New("expectedProfit == 0 or 1")
		}
		tf.Trade.TradeMethod = swapExactOutputSingle
		tf.InitialPairV3 = pd.V3Pair.ConvertToJSONType()
		a.GetMetricsClient().StageProgressionMetrics.CountPostProcessTx(float64(1))
		a.GetMetricsClient().TxFetcherMetrics.TransactionGroup(toAddr, swapExactOutputSingle)
		a.GetMetricsClient().TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.TokenFeePath.TokenIn.String(), inputs.TokenFeePath.GetEndToken().String())
		a.GetMetricsClient().TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactOutputSingle, pd.V3Pair.PoolAddress, tf.FrontRunTrade.AmountInAddr.String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		tfSlice = append(tfSlice, tf)
		log.Info().Msg("saving mempool tx")
		err = a.SaveMempoolTx(ctx, bn, []web3_client.TradeExecutionFlowJSON{tf})
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
	case swapExactTokensForTokens:
		inputs := &web3_client.SwapExactTokensForTokensParamsV3{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			//log.Err(err).Msg("swapExactTokensForTokens: failed to decode swap exact tokens for tokens args")
			return nil, err
		}
		pd, err := a.GetUniswapClient().GetV2PricingData(ctx, inputs.Path)
		if err != nil {
			if pd != nil {
				a.GetMetricsClient().ErrTrackingMetrics.RecordError(swapExactTokensForTokens, pd.V2Pair.PairContractAddr)
			}
			//log.Err(err).Msg("failed to get pricing data")
			return nil, err
		}
		if pd == nil {
			return nil, errors.New("pd is nil")
		}
		tf := inputs.BinarySearch(pd.V2Pair)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" || tf.SandwichPrediction.ExpectedProfit == "" {
			return nil, errors.New("expectedProfit == 0 or 1")
		}
		newTx := artemis_trading_types.JSONTx{}
		err = newTx.UnmarshalTx(tx.Tx)
		if err != nil {
			log.Err(err).Msg("failed to unmarshal tx")
			return nil, err
		}
		tf.Tx = newTx
		tf.Trade.TradeMethod = swapExactTokensForTokens
		tf.InitialPair = pd.V2Pair.ConvertToJSONType()
		a.GetMetricsClient().StageProgressionMetrics.CountPostProcessTx(float64(1))
		a.GetMetricsClient().TxFetcherMetrics.TransactionGroup(toAddr, swapExactTokensForTokens)
		pend := len(inputs.Path) - 1
		a.GetMetricsClient().TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.Path[0].String(), inputs.Path[pend].String())
		a.GetMetricsClient().TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactTokensForTokens, pd.V2Pair.PairContractAddr, inputs.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		tfSlice = append(tfSlice, tf)
		log.Info().Msg("saving mempool tx")
		err = a.SaveMempoolTx(ctx, bn, []web3_client.TradeExecutionFlowJSON{tf})
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
	case swapExactInputMultihop:
		a.GetMetricsClient().TxFetcherMetrics.TransactionGroup(toAddr, swapExactInputMultihop)
	case swapExactOutputMultihop:
		a.GetMetricsClient().TxFetcherMetrics.TransactionGroup(toAddr, swapExactOutputMultihop)
	}
	return tfSlice, nil
}
