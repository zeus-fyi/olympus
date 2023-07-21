package artemis_realtime_trading

import (
	"context"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	metrics_trading "github.com/zeus-fyi/olympus/pkg/apollo/ethereum/mev/trading"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
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

func RealTimeProcessUniswapV3RouterTx(ctx context.Context, tx web3_client.MevTx, abiFile *abi.ABI, filter *strings_filter.FilterOpts, m *metrics_trading.TradingMetrics, w3a web3_actions.Web3Actions) ([]web3_client.TradeExecutionFlow, error) {
	toAddr := tx.Tx.To().String()
	bn, berr := artemis_trading_cache.GetLatestBlock(context.Background())
	if berr != nil {
		log.Err(berr).Msg("RealTimeProcessUniswapV3RouterTx: failed to get latest block")
		return nil, errors.New("failed to get latest block")
	}
	var tfSlice []web3_client.TradeExecutionFlow
	if strings.HasPrefix(tx.MethodName, multicall) {
		m.TxFetcherMetrics.TransactionGroup(toAddr, multicall)
		inputs := &web3_client.Multicall{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			log.Err(err).Msg("RealTimeProcessUniswapV3RouterTx: failed to decode multicall args")
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
			tf, terr := processUniswapV3Txs(ctx, bn, newTx, m, w3a)
			if terr != nil {
				log.Err(terr).Msg("failed to process uniswap v3 tx")
				continue
			}
			tfSlice = append(tfSlice, tf...)
		}
	} else {
		tf, err := processUniswapV3Txs(ctx, bn, tx, m, w3a)
		if err != nil {
			log.Err(err).Msg("failed to process uniswap v3 tx")
			return nil, err
		}
		tfSlice = append(tfSlice, tf...)
	}
	if len(tfSlice) == 0 {
		return nil, errors.New("RealTimeProcessUniswapV3RouterTx: tfSlice is empty")
	}
	return tfSlice, nil
}

func processUniswapV3Txs(ctx context.Context, bn uint64, tx web3_client.MevTx, m *metrics_trading.TradingMetrics, w3a web3_actions.Web3Actions) ([]web3_client.TradeExecutionFlow, error) {
	var tfSlice []web3_client.TradeExecutionFlow
	toAddr := tx.Tx.To().String()
	switch tx.MethodName {
	case exactInput:
		inputs := &web3_client.ExactInputParams{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			log.Err(err).Msg("failed to decode exact input args")
			return nil, err
		}
		pd, err := uniswap_pricing.GetV3PricingData(ctx, w3a, inputs.TokenFeePath)
		if err != nil {
			if pd != nil && m != nil {
				m.ErrTrackingMetrics.RecordError(exactInput, pd.V3Pair.PoolAddress)
			}
			//log.Err(err).Msg("failed to get pricing data")
			return nil, err
		}
		if pd == nil {
			return nil, errors.New("pd is nil")
		}
		tf, err := inputs.BinarySearch(pd)
		if err != nil {
			return nil, err
		}
		tf.CurrentBlockNumber = artemis_eth_units.NewBigIntFromUint(bn)
		tf.Tx = tx.Tx
		err = ApplyMaxTransferTax(ctx, &tf)
		if err != nil {
			return nil, err
		}
		log.Info().Msg("exactInput: saving mempool tx")
		err = SaveMempoolTx(ctx, []web3_client.TradeExecutionFlow{tf}, m)
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
		if m != nil {
			m.TxFetcherMetrics.TransactionGroup(toAddr, exactInput)
			m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.TokenFeePath.TokenIn.String(), inputs.TokenFeePath.GetEndToken().String())
			m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, exactInput, pd.V3Pair.PoolAddress, inputs.TokenFeePath.TokenIn.String(), tf.SandwichPrediction.SellAmount.String(), tf.SandwichPrediction.ExpectedProfit.String())
			m.StageProgressionMetrics.CountPostProcessTx(float64(1))
		}
		tfSlice = append(tfSlice, tf)
	case exactOutput:
		inputs := &web3_client.ExactOutputParams{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			log.Err(err).Msg("failed to decode exact output args")
			return nil, err
		}
		pd, err := uniswap_pricing.GetV3PricingData(ctx, w3a, inputs.TokenFeePath)
		if err != nil {
			if pd != nil && m != nil {
				m.ErrTrackingMetrics.RecordError(exactOutput, pd.V3Pair.PoolAddress)
			}
			//log.Err(err).Msg("failed to get pricing data")
			return nil, err
		}
		if pd == nil {
			return nil, errors.New("pd is nil")
		}
		tf, err := inputs.BinarySearch(pd)
		if err != nil {
			return nil, err
		}
		tf.CurrentBlockNumber = artemis_eth_units.NewBigIntFromUint(bn)
		tf.Tx = tx.Tx
		err = ApplyMaxTransferTax(ctx, &tf)
		if err != nil {
			return nil, err
		}
		log.Info().Msg("exactOutput: saving mempool tx")
		err = SaveMempoolTx(ctx, []web3_client.TradeExecutionFlow{tf}, m)
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
		if m != nil {
			m.StageProgressionMetrics.CountPostProcessTx(float64(1))
			m.TxFetcherMetrics.TransactionGroup(toAddr, exactOutput)
			m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.TokenFeePath.TokenIn.String(), inputs.TokenFeePath.GetEndToken().String())
			m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, exactOutput, pd.V3Pair.PoolAddress, tf.FrontRunTrade.AmountInAddr.String(), tf.SandwichPrediction.SellAmount.String(), tf.SandwichPrediction.ExpectedProfit.String())
		}
		tfSlice = append(tfSlice, tf)
	case swapExactInputSingle:
		inputs := &web3_client.SwapExactInputSingleArgs{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			log.Err(err).Msg("failed to decode swap exact input single args")
			return nil, err
		}
		pd, err := uniswap_pricing.GetV3PricingData(ctx, w3a, inputs.TokenFeePath)
		if err != nil {
			if pd != nil && m != nil {
				m.ErrTrackingMetrics.RecordError(swapExactInputSingle, pd.V3Pair.PoolAddress)
			}
			//log.Err(err).Msg("failed to get pricing data")
			return nil, err
		}
		if pd == nil {
			return nil, errors.New("pd is nil")
		}
		tf, err := inputs.BinarySearch(pd)
		if err != nil {
			return nil, err
		}
		tf.CurrentBlockNumber = artemis_eth_units.NewBigIntFromUint(bn)
		tf.Tx = tx.Tx
		err = ApplyMaxTransferTax(ctx, &tf)
		if err != nil {
			return nil, err
		}
		log.Info().Msg("saving mempool tx")
		err = SaveMempoolTx(ctx, []web3_client.TradeExecutionFlow{tf}, m)
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
		if m != nil {
			m.StageProgressionMetrics.CountPostProcessTx(float64(1))
			m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactInputSingle)
			m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.TokenFeePath.TokenIn.String(), inputs.TokenFeePath.GetEndToken().String())
			m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactInputSingle, pd.V3Pair.PoolAddress, inputs.TokenFeePath.TokenIn.String(), tf.SandwichPrediction.SellAmount.String(), tf.SandwichPrediction.ExpectedProfit.String())
		}
		tfSlice = append(tfSlice, tf)
	case swapExactOutputSingle:
		inputs := &web3_client.SwapExactOutputSingleArgs{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			//log.Err(err).Msg("failed to decode swap exact output single args")
			return nil, err
		}
		pd, err := uniswap_pricing.GetV3PricingData(ctx, w3a, inputs.TokenFeePath)
		if err != nil {
			if pd != nil && m != nil {
				m.ErrTrackingMetrics.RecordError(swapExactOutputSingle, pd.V3Pair.PoolAddress)
			}
			//log.Err(err).Msg("failed to get pricing data")
			return nil, err
		}
		if pd == nil {
			return nil, errors.New("pd is nil")
		}
		tf, err := inputs.BinarySearch(pd)
		if err != nil {
			return nil, err
		}
		tf.CurrentBlockNumber = artemis_eth_units.NewBigIntFromUint(bn)
		tf.Tx = tx.Tx
		err = ApplyMaxTransferTax(ctx, &tf)
		if err != nil {
			return nil, err
		}
		log.Info().Msg("saving mempool tx")
		err = SaveMempoolTx(ctx, []web3_client.TradeExecutionFlow{tf}, m)
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
		if m != nil {
			m.StageProgressionMetrics.CountPostProcessTx(float64(1))
			m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactOutputSingle)
			m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.TokenFeePath.TokenIn.String(), inputs.TokenFeePath.GetEndToken().String())
			m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactOutputSingle, pd.V3Pair.PoolAddress, tf.FrontRunTrade.AmountInAddr.String(), tf.SandwichPrediction.SellAmount.String(), tf.SandwichPrediction.ExpectedProfit.String())
		}
		tfSlice = append(tfSlice, tf)
	case swapExactTokensForTokens:
		inputs := &web3_client.SwapExactTokensForTokensParamsV3{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			//log.Err(err).Msg("swapExactTokensForTokens: failed to decode swap exact tokens for tokens args")
			return nil, err
		}
		pd, err := uniswap_pricing.GetV2PricingData(ctx, w3a, inputs.Path)
		if err != nil {
			if pd != nil && m != nil {
				m.ErrTrackingMetrics.RecordError(swapExactTokensForTokens, pd.V2Pair.PairContractAddr)
			}
			//log.Err(err).Msg("failed to get pricing data")
			return nil, err
		}
		if pd == nil {
			return nil, errors.New("pd is nil")
		}
		tf, err := inputs.BinarySearch(pd.V2Pair)
		if err != nil {
			return nil, err
		}
		tf.CurrentBlockNumber = artemis_eth_units.NewBigIntFromUint(bn)
		tf.Tx = tx.Tx
		err = ApplyMaxTransferTax(ctx, &tf)
		if err != nil {
			return nil, err
		}
		log.Info().Msg("saving mempool tx")
		err = SaveMempoolTx(ctx, []web3_client.TradeExecutionFlow{tf}, m)
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
		if m != nil {
			m.StageProgressionMetrics.CountPostProcessTx(float64(1))
			m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactTokensForTokens)
			pend := len(inputs.Path) - 1
			m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.Path[0].String(), inputs.Path[pend].String())
			m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactTokensForTokens, pd.V2Pair.PairContractAddr, inputs.Path[0].String(), tf.SandwichPrediction.SellAmount.String(), tf.SandwichPrediction.ExpectedProfit.String())
		}
		tfSlice = append(tfSlice, tf)
	case swapExactInputMultihop:
		if m != nil {
			m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactInputMultihop)
		}
	case swapExactOutputMultihop:
		if m != nil {
			m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactOutputMultihop)
		}
	}
	return tfSlice, nil
}
