package artemis_realtime_trading

import (
	"context"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/rs/zerolog/log"
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

func (a *ActiveTrading) RealTimeProcessUniswapV3RouterTx(ctx context.Context, tx web3_client.MevTx, abiFile *abi.ABI, filter *strings_filter.FilterOpts) ([]*web3_client.TradeExecutionFlowJSON, error) {
	toAddr := tx.Tx.To().String()
	var tfSlice []*web3_client.TradeExecutionFlowJSON
	if strings.HasPrefix(tx.MethodName, multicall) {
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, multicall)
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
				return nil, terr
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

func (a *ActiveTrading) processUniswapV3Txs(ctx context.Context, tx web3_client.MevTx) ([]*web3_client.TradeExecutionFlowJSON, error) {
	var tfSlice []*web3_client.TradeExecutionFlowJSON
	toAddr := tx.Tx.To().String()
	switch tx.MethodName {
	case exactInput:
		inputs := &web3_client.ExactInputParams{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			log.Err(err).Msg("failed to decode exact input args")
			return nil, err
		}
		pd, err := a.u.GetV3PricingData(ctx, inputs.TokenFeePath)
		if err != nil {
			a.m.ErrTrackingMetrics.RecordError(exactInput, pd.V3Pair.PoolAddress)
			log.Err(err).Msg("failed to get pricing data")
			return nil, err
		}
		tf := inputs.BinarySearch(pd)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" {
			return nil, errors.New("expectedProfit == 0 or 1")
		}
		tf.InitialPairV3 = pd.V3Pair.ConvertToJSONType()
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, exactInput)
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.TokenFeePath.TokenIn.String(), inputs.TokenFeePath.GetEndToken().String())
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, exactInput, pd.V3Pair.PoolAddress, inputs.TokenFeePath.TokenIn.String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		tfSlice = append(tfSlice, &tf)
	case exactOutput:
		inputs := &web3_client.ExactOutputParams{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			log.Err(err).Msg("failed to decode exact output args")
			return nil, err
		}
		pd, err := a.u.GetV3PricingData(ctx, inputs.TokenFeePath)
		if err != nil {
			a.m.ErrTrackingMetrics.RecordError(exactOutput, pd.V3Pair.PoolAddress)
			log.Err(err).Msg("failed to get pricing data")
			return nil, err
		}
		tf := inputs.BinarySearch(pd)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" {
			return nil, errors.New("expectedProfit == 0 or 1")
		}
		tf.InitialPairV3 = pd.V3Pair.ConvertToJSONType()
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, exactOutput)
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.TokenFeePath.TokenIn.String(), inputs.TokenFeePath.GetEndToken().String())
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, exactOutput, pd.V3Pair.PoolAddress, tf.FrontRunTrade.AmountInAddr.String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		tfSlice = append(tfSlice, &tf)
	case swapExactInputSingle:
		inputs := &web3_client.SwapExactInputSingleArgs{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			log.Err(err).Msg("failed to decode swap exact input single args")
			return nil, err
		}
		pd, err := a.u.GetV3PricingData(ctx, inputs.TokenFeePath)
		if err != nil {
			a.m.ErrTrackingMetrics.RecordError(swapExactInputSingle, pd.V3Pair.PoolAddress)
			log.Err(err).Msg("failed to get pricing data")
			return nil, err
		}
		tf := inputs.BinarySearch(pd)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" {
			return nil, errors.New("expectedProfit == 0 or 1")
		}
		tf.InitialPairV3 = pd.V3Pair.ConvertToJSONType()
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactInputSingle)
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.TokenFeePath.TokenIn.String(), inputs.TokenFeePath.GetEndToken().String())
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactInputSingle, pd.V3Pair.PoolAddress, inputs.TokenFeePath.TokenIn.String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		tfSlice = append(tfSlice, &tf)
	case swapExactOutputSingle:
		inputs := &web3_client.SwapExactOutputSingleArgs{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			log.Err(err).Msg("failed to decode swap exact output single args")
			return nil, err
		}
		pd, err := a.u.GetV3PricingData(ctx, inputs.TokenFeePath)
		if err != nil {
			a.m.ErrTrackingMetrics.RecordError(swapExactOutputSingle, pd.V3Pair.PoolAddress)
			log.Err(err).Msg("failed to get pricing data")
			return nil, err
		}
		tf := inputs.BinarySearch(pd)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" {
			return nil, errors.New("expectedProfit == 0 or 1")
		}
		tf.InitialPairV3 = pd.V3Pair.ConvertToJSONType()
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactOutputSingle)
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.TokenFeePath.TokenIn.String(), inputs.TokenFeePath.GetEndToken().String())
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactOutputSingle, pd.V3Pair.PoolAddress, tf.FrontRunTrade.AmountInAddr.String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		tfSlice = append(tfSlice, &tf)
	case swapExactTokensForTokens:
		inputs := &web3_client.SwapExactTokensForTokensParamsV3{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			log.Err(err).Msg("swapExactTokensForTokens: failed to decode swap exact tokens for tokens args")
			return nil, err
		}
		pd, err := a.u.GetV2PricingData(ctx, inputs.Path)
		if err != nil {
			a.m.ErrTrackingMetrics.RecordError(swapExactTokensForTokens, pd.V2Pair.PairContractAddr)
			log.Err(err).Msg("failed to get pricing data")
			return nil, err
		}
		tf := inputs.BinarySearch(pd.V2Pair)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" {
			return nil, errors.New("expectedProfit == 0 or 1")
		}
		tf.InitialPair = pd.V2Pair.ConvertToJSONType()
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactTokensForTokens)
		pend := len(inputs.Path) - 1
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.Path[0].String(), inputs.Path[pend].String())
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactTokensForTokens, pd.V2Pair.PairContractAddr, inputs.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		tfSlice = append(tfSlice, &tf)
	case swapExactInputMultihop:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactInputMultihop)
	case swapExactOutputMultihop:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactOutputMultihop)
	}
	return tfSlice, nil
}
