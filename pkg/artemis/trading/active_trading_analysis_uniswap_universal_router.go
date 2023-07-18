package artemis_realtime_trading

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	metrics_trading "github.com/zeus-fyi/olympus/pkg/apollo/ethereum/mev/trading"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func RealTimeProcessUniversalRouterTx(ctx context.Context, tx web3_client.MevTx, m *metrics_trading.TradingMetrics, w3a web3_actions.Web3Actions) ([]web3_client.TradeExecutionFlowJSON, error) {
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
			pd, perr := uniswap_pricing.GetV3PricingData(ctx, w3a, inputs.Path)
			if perr != nil {
				if pd != nil && m != nil {
					m.ErrTrackingMetrics.RecordError(web3_client.V3SwapExactIn, pd.V3Pair.PoolAddress)
				}
				//log.Err(perr).Msg("V3SwapExactIn: error getting pricing data")
				return nil, perr
			}
			if pd == nil {
				return nil, errors.New("pd is nil")
			}
			tf := inputs.BinarySearch(pd)
			tf.Tx.Hash = tx.Tx.Hash().String()
			err = ApplyMaxTransferTax(&tf)
			if err != nil {
				return nil, err
			}
			newJsonTx := artemis_trading_types.JSONTx{}
			err = newJsonTx.UnmarshalTx(tx.Tx)
			if err != nil {
				return nil, err
			}
			if artemis_eth_units.IsStrXLessThanEqZeroOrOne(tf.SandwichPrediction.ExpectedProfit) || artemis_eth_units.IsStrXLessThanEqZeroOrOne(tf.SandwichTrade.AmountOut) {
				return nil, errors.New("expectedProfit == 0 or 1")
			}
			tf.Tx = newJsonTx
			tf.InitialPairV3 = pd.V3Pair.ConvertToJSONType()
			tf.Trade.TradeMethod = web3_client.V3SwapExactIn

			if m != nil {
				m.StageProgressionMetrics.CountPostProcessTx(float64(1))
				m.TxFetcherMetrics.TransactionGroup(toAddr, web3_client.V3SwapExactIn)
				m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.Path.TokenIn.String(), inputs.Path.GetEndToken().String())
				m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, web3_client.V3SwapExactIn, pd.V3Pair.PoolAddress, inputs.Path.TokenIn.String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
			}

			tfSlice = append(tfSlice, tf)
			log.Info().Msg("saving mempool tx")
			err = SaveMempoolTx(ctx, bn, []web3_client.TradeExecutionFlowJSON{tf}, m)
			if err != nil {
				log.Err(err).Msg("failed to save mempool tx")
				return nil, errors.New("failed to save mempool tx")
			}
		case web3_client.V3SwapExactOut:
			//fmt.Println("V3SwapExactOut: ProcessUniversalRouterTxs")
			inputs := subtx.DecodedInputs.(web3_client.V3SwapExactOutParams)
			pd, perr := uniswap_pricing.GetV3PricingData(ctx, w3a, inputs.Path)
			if perr != nil {
				if pd != nil && m != nil {
					m.ErrTrackingMetrics.RecordError(web3_client.V3SwapExactOut, pd.V3Pair.PoolAddress)
				}
				//log.Err(perr).Msg("V3SwapExactIn: error getting pricing data")
				return nil, perr
			}
			if pd == nil {
				return nil, errors.New("pd is nil")
			}
			tf := inputs.BinarySearch(pd)
			tf.Tx.Hash = tx.Tx.Hash().String()
			err = ApplyMaxTransferTax(&tf)
			if err != nil {
				return nil, err
			}
			newJsonTx := artemis_trading_types.JSONTx{}
			err = newJsonTx.UnmarshalTx(tx.Tx)
			if err != nil {
				return nil, err
			}
			tf.Tx = newJsonTx
			tf.InitialPairV3 = pd.V3Pair.ConvertToJSONType()
			tf.Trade.TradeMethod = web3_client.V3SwapExactOut
			if m != nil {
				m.StageProgressionMetrics.CountPostProcessTx(float64(1))
				m.TxFetcherMetrics.TransactionGroup(toAddr, web3_client.V3SwapExactOut)
				m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.Path.TokenIn.String(), inputs.Path.GetEndToken().String())
				m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, web3_client.V3SwapExactOut, pd.V3Pair.PoolAddress, tf.FrontRunTrade.AmountInAddr.String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
			}
			tfSlice = append(tfSlice, tf)
			log.Info().Msg("saving mempool tx")
			err = SaveMempoolTx(ctx, bn, []web3_client.TradeExecutionFlowJSON{tf}, m)
			if err != nil {
				log.Err(err).Msg("failed to save mempool tx")
				return nil, errors.New("failed to save mempool tx")
			}
		case web3_client.V2SwapExactIn:
			//fmt.Println("V2SwapExactIn: ProcessUniversalRouterTxs")
			inputs := subtx.DecodedInputs.(web3_client.V2SwapExactInParams)
			pd, perr := uniswap_pricing.GetV2PricingData(ctx, w3a, inputs.Path)
			if perr != nil {
				if pd != nil && m != nil {
					m.ErrTrackingMetrics.RecordError(web3_client.V2SwapExactIn, pd.V2Pair.PairContractAddr)
				}
				//log.Err(perr).Msg("V2SwapExactIn: error getting pricing data")
				return nil, perr
			}
			if pd == nil {
				return nil, errors.New("pd is nil")
			}
			tf := inputs.BinarySearch(pd.V2Pair)
			tf.Tx.Hash = tx.Tx.Hash().String()
			err = ApplyMaxTransferTax(&tf)
			if err != nil {
				return nil, err
			}

			newJsonTx := artemis_trading_types.JSONTx{}
			err = newJsonTx.UnmarshalTx(tx.Tx)
			if err != nil {
				return nil, err
			}
			tf.Tx = newJsonTx
			tf.Trade.TradeMethod = web3_client.V2SwapExactIn
			tf.InitialPair = pd.V2Pair.ConvertToJSONType()
			if m != nil {
				m.StageProgressionMetrics.CountPostProcessTx(float64(1))
				m.TxFetcherMetrics.TransactionGroup(toAddr, web3_client.V2SwapExactIn)
				pend := len(inputs.Path) - 1
				m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.Path[0].String(), inputs.Path[pend].String())
				m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, web3_client.V2SwapExactIn, pd.V2Pair.PairContractAddr, inputs.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
			}
			tfSlice = append(tfSlice, tf)
			log.Info().Msg("saving mempool tx")
			err = SaveMempoolTx(ctx, bn, []web3_client.TradeExecutionFlowJSON{tf}, m)
			if err != nil {
				log.Err(err).Msg("failed to save mempool tx")
				return nil, errors.New("failed to save mempool tx")
			}
		case web3_client.V2SwapExactOut:
			//fmt.Println("V2SwapExactOut: ProcessUniversalRouterTxs")
			inputs := subtx.DecodedInputs.(web3_client.V2SwapExactOutParams)
			pd, perr := uniswap_pricing.GetV2PricingData(ctx, w3a, inputs.Path)
			if perr != nil {
				if pd != nil && m != nil {
					m.ErrTrackingMetrics.RecordError(web3_client.V2SwapExactOut, pd.V2Pair.PairContractAddr)
				}
				//log.Err(perr).Msg("V2SwapExactOut: error getting pricing data")
				return nil, perr
			}
			if pd == nil {
				return nil, errors.New("pd is nil")
			}
			tf := inputs.BinarySearch(pd.V2Pair)
			tf.Tx.Hash = tx.Tx.Hash().String()
			err = ApplyMaxTransferTax(&tf)
			if err != nil {
				return nil, err
			}
			newJsonTx := artemis_trading_types.JSONTx{}
			err = newJsonTx.UnmarshalTx(tx.Tx)
			if err != nil {
				return nil, err
			}
			tf.Tx = newJsonTx
			tf.Trade.TradeMethod = web3_client.V2SwapExactOut
			tf.InitialPair = pd.V2Pair.ConvertToJSONType()
			pend := len(inputs.Path) - 1
			if m != nil {
				m.StageProgressionMetrics.CountPostProcessTx(float64(1))
				m.TxFetcherMetrics.TransactionGroup(toAddr, web3_client.V2SwapExactOut)
				m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.Path[0].String(), inputs.Path[pend].String())
				m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, web3_client.V2SwapExactOut, pd.V2Pair.PairContractAddr, tf.FrontRunTrade.AmountInAddr.String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
			}
			tfSlice = append(tfSlice, tf)
			log.Info().Msg("saving mempool tx")
			err = SaveMempoolTx(ctx, bn, []web3_client.TradeExecutionFlowJSON{tf}, m)
			if err != nil {
				log.Err(err).Msg("failed to save mempool tx")
				return nil, errors.New("failed to save mempool tx")
			}
		default:
		}
	}
	return tfSlice, nil
}
