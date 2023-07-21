package artemis_realtime_trading

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	metrics_trading "github.com/zeus-fyi/olympus/pkg/apollo/ethereum/mev/trading"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func RealTimeProcessUniversalRouterTx(ctx context.Context, tx web3_client.MevTx, m *metrics_trading.TradingMetrics, w3a web3_actions.Web3Actions, abiFile *abi.ABI) ([]web3_client.TradeExecutionFlow, error) {
	if tx.Tx == nil || tx.Args == nil || tx.Tx.To() == nil {
		return nil, errors.New("tx is nil")
	}
	subcmd, err := web3_client.NewDecodedUniversalRouterExecCmdFromMap(tx.Args, abiFile)
	if err != nil {
		log.Err(err).Msg("RealTimeProcessUniversalRouterTx: failed to decode universal router exec cmd")
		return nil, err
	}
	bn, err := artemis_trading_cache.GetLatestBlock(context.Background())
	if err != nil {
		log.Err(err).Msg("RealTimeProcessUniversalRouterTx: failed to get latest block")
		return nil, errors.New("failed to get latest block")
	}
	var tfSlice []web3_client.TradeExecutionFlow
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
				continue
			}
			if pd == nil {
				continue
			}
			tf, terr := inputs.BinarySearch(pd)
			if terr != nil {
				continue
			}
			tf.CurrentBlockNumber = artemis_eth_units.NewBigIntFromUint(bn)
			tf.Tx = tx.Tx
			err = ApplyMaxTransferTax(ctx, &tf)
			if err != nil {
				continue
			}
			log.Info().Msg("V3SwapExactIn: saving mempool tx")
			err = SaveMempoolTx(ctx, []web3_client.TradeExecutionFlow{tf}, m)
			if err != nil {
				log.Err(err).Msg("failed to save mempool tx")
				continue
			}
			if m != nil {
				m.StageProgressionMetrics.CountPostProcessTx(float64(1))
				m.TxFetcherMetrics.TransactionGroup(toAddr, web3_client.V3SwapExactIn)
				m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.Path.TokenIn.String(), inputs.Path.GetEndToken().String())
				m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, web3_client.V3SwapExactIn, pd.V3Pair.PoolAddress, inputs.Path.TokenIn.String(), tf.SandwichPrediction.SellAmount.String(), tf.SandwichPrediction.ExpectedProfit.String())
			}
			tfSlice = append(tfSlice, tf)
		case web3_client.V3SwapExactOut:
			//fmt.Println("V3SwapExactOut: ProcessUniversalRouterTxs")
			inputs := subtx.DecodedInputs.(web3_client.V3SwapExactOutParams)
			pd, perr := uniswap_pricing.GetV3PricingData(ctx, w3a, inputs.Path)
			if perr != nil {
				if pd != nil && m != nil {
					m.ErrTrackingMetrics.RecordError(web3_client.V3SwapExactOut, pd.V3Pair.PoolAddress)
				}
				log.Err(perr).Msg("V3SwapExactIn: error getting pricing data")
				continue
			}
			if pd == nil {
				continue
			}
			tf, terr := inputs.BinarySearch(pd)
			if terr != nil {
				continue
			}
			tf.CurrentBlockNumber = artemis_eth_units.NewBigIntFromUint(bn)
			tf.Tx = tx.Tx
			err = ApplyMaxTransferTax(ctx, &tf)
			if err != nil {
				continue
			}
			log.Info().Msg("V3SwapExactOut: saving mempool tx")
			err = SaveMempoolTx(ctx, []web3_client.TradeExecutionFlow{tf}, m)
			if err != nil {
				log.Err(err).Msg("failed to save mempool tx")
				continue
			}
			if m != nil {
				m.StageProgressionMetrics.CountPostProcessTx(float64(1))
				m.TxFetcherMetrics.TransactionGroup(toAddr, web3_client.V3SwapExactOut)
				m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.Path.TokenIn.String(), inputs.Path.GetEndToken().String())
				m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, web3_client.V3SwapExactOut, pd.V3Pair.PoolAddress, tf.FrontRunTrade.AmountInAddr.String(), tf.SandwichPrediction.SellAmount.String(), tf.SandwichPrediction.ExpectedProfit.String())
			}
			tfSlice = append(tfSlice, tf)
		case web3_client.V2SwapExactIn:
			//fmt.Println("V2SwapExactIn: ProcessUniversalRouterTxs")
			inputs := subtx.DecodedInputs.(web3_client.V2SwapExactInParams)
			pd, perr := uniswap_pricing.GetV2PricingData(ctx, w3a, inputs.Path)
			if perr != nil {
				if pd != nil && m != nil {
					m.ErrTrackingMetrics.RecordError(web3_client.V2SwapExactIn, pd.V2Pair.PairContractAddr)
				}
				log.Err(perr).Msg("V2SwapExactIn: error getting pricing data")
				continue
			}
			if pd == nil {
				continue
			}
			tf, terr := inputs.BinarySearch(pd.V2Pair)
			if terr != nil {
				continue
			}
			tf.CurrentBlockNumber = artemis_eth_units.NewBigIntFromUint(bn)
			tf.Tx = tx.Tx
			err = ApplyMaxTransferTax(ctx, &tf)
			if err != nil {
				continue
			}
			if m != nil {
				m.StageProgressionMetrics.CountPostProcessTx(float64(1))
				m.TxFetcherMetrics.TransactionGroup(toAddr, web3_client.V2SwapExactIn)
				pend := len(inputs.Path) - 1
				m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.Path[0].String(), inputs.Path[pend].String())
				m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, web3_client.V2SwapExactIn, pd.V2Pair.PairContractAddr, inputs.Path[0].String(), tf.SandwichPrediction.SellAmount.String(), tf.SandwichPrediction.ExpectedProfit.String())
			}
			log.Info().Msg("V2SwapExactIn: saving mempool tx")
			err = SaveMempoolTx(ctx, []web3_client.TradeExecutionFlow{tf}, m)
			if err != nil {
				continue
			}
			tfSlice = append(tfSlice, tf)
		case web3_client.V2SwapExactOut:
			//fmt.Println("V2SwapExactOut: ProcessUniversalRouterTxs")
			inputs := subtx.DecodedInputs.(web3_client.V2SwapExactOutParams)
			pd, perr := uniswap_pricing.GetV2PricingData(ctx, w3a, inputs.Path)
			if perr != nil {
				if pd != nil && m != nil {
					m.ErrTrackingMetrics.RecordError(web3_client.V2SwapExactOut, pd.V2Pair.PairContractAddr)
				}
				//log.Err(perr).Msg("V2SwapExactOut: error getting pricing data")
				continue
			}
			if pd == nil {
				continue
			}
			tf, terr := inputs.BinarySearch(pd.V2Pair)
			if terr != nil {
				continue
			}
			tf.CurrentBlockNumber = artemis_eth_units.NewBigIntFromUint(bn)
			tf.Tx = tx.Tx
			err = ApplyMaxTransferTax(ctx, &tf)
			if err != nil {
				continue
			}
			pend := len(inputs.Path) - 1
			if m != nil {
				m.StageProgressionMetrics.CountPostProcessTx(float64(1))
				m.TxFetcherMetrics.TransactionGroup(toAddr, web3_client.V2SwapExactOut)
				m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, inputs.Path[0].String(), inputs.Path[pend].String())
				m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, web3_client.V2SwapExactOut, pd.V2Pair.PairContractAddr, tf.FrontRunTrade.AmountInAddr.String(), tf.SandwichPrediction.SellAmount.String(), tf.SandwichPrediction.ExpectedProfit.String())
			}
			log.Info().Msg("V2SwapExactOut: saving mempool tx")
			err = SaveMempoolTx(ctx, []web3_client.TradeExecutionFlow{tf}, m)
			if err != nil {
				continue
			}
			tfSlice = append(tfSlice, tf)
		default:
		}
	}
	if len(tfSlice) == 0 {
		return nil, errors.New("RealTimeProcessUniversalRouterTx: no txs to process")
	}
	return tfSlice, nil
}
