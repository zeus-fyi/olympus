package artemis_realtime_trading

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	metrics_trading "github.com/zeus-fyi/olympus/pkg/apollo/ethereum/mev/trading"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

const (
	addLiquidity                 = "addLiquidity"
	addLiquidityETH              = "addLiquidityETH"
	removeLiquidity              = "removeLiquidity"
	removeLiquidityETH           = "removeLiquidityETH"
	removeLiquidityWithPermit    = "removeLiquidityWithPermit"
	removeLiquidityETHWithPermit = "removeLiquidityETHWithPermit"
	swapExactTokensForTokens     = "swapExactTokensForTokens"
	swapTokensForExactTokens     = "swapTokensForExactTokens"
	swapExactETHForTokens        = "swapExactETHForTokens"
	swapTokensForExactETH        = "swapTokensForExactETH"
	swapExactTokensForETH        = "swapExactTokensForETH"
	swapETHForExactTokens        = "swapETHForExactTokens"

	// UniswapV2Router02
	swapExactTokensForETHSupportingFeeOnTransferTokens    = "swapExactTokensForETHSupportingFeeOnTransferTokens"
	swapExactETHForTokensSupportingFeeOnTransferTokens    = "swapExactETHForTokensSupportingFeeOnTransferTokens"
	swapExactTokensForTokensSupportingFeeOnTransferTokens = "swapExactTokensForTokensSupportingFeeOnTransferTokens"

	removeLiquidityETHWithPermitSupportingFeeOnTransferTokens = "removeLiquidityETHWithPermitSupportingFeeOnTransferTokens"
	removeLiquidityETHSupportingFeeOnTransferTokens           = "removeLiquidityETHSupportingFeeOnTransferTokens"
)

func RealTimeProcessUniswapV2RouterTx(ctx context.Context, tx web3_client.MevTx, m *metrics_trading.TradingMetrics, w3a web3_actions.Web3Actions, abiFile *abi.ABI) ([]web3_client.TradeExecutionFlow, error) {
	bn, berr := artemis_trading_cache.GetLatestBlock(ctx)
	if berr != nil {
		log.Err(berr).Msg("RealTimeProcessUniswapV2RouterTx: failed to get latest block")
		return nil, errors.New("RealTimeProcessUniswapV2RouterTx: failed to get latest block")
	}
	toAddr := tx.Tx.To().String()
	var tfSlice []web3_client.TradeExecutionFlow
	switch tx.MethodName {
	case addLiquidity:
		if m != nil {
			m.TxFetcherMetrics.TransactionGroup(toAddr, addLiquidity)
		}
	case addLiquidityETH:
		if tx.Tx.Value() == nil {
			return nil, errors.New("RealTimeProcessUniswapV2RouterTx: addLiquidityETH tx has no value")
		}
		if m != nil {
			m.TxFetcherMetrics.TransactionGroup(toAddr, addLiquidityETH)
		}
	case removeLiquidity:
		if m != nil {
			m.TxFetcherMetrics.TransactionGroup(toAddr, removeLiquidity)
		}
	case removeLiquidityETH:
		if m != nil {
			m.TxFetcherMetrics.TransactionGroup(toAddr, removeLiquidityETH)
		}
	case removeLiquidityWithPermit:
		if m != nil {
			m.TxFetcherMetrics.TransactionGroup(toAddr, removeLiquidityWithPermit)
		}
	case removeLiquidityETHWithPermit:
		if m != nil {
			m.TxFetcherMetrics.TransactionGroup(toAddr, removeLiquidityETHWithPermit)
		}
	case swapExactTokensForTokens:
		st := web3_client.SwapExactTokensForTokensParams{}
		err := st.Decode(ctx, tx.Args)
		if err != nil {
			return nil, err
		}
		pd, err := uniswap_pricing.GetV2PricingData(ctx, w3a, st.Path)
		if err != nil {
			if pd != nil && m != nil {
				m.ErrTrackingMetrics.RecordError(swapExactTokensForTokens, pd.V2Pair.PairContractAddr)
				return nil, err
			}
		}
		if pd == nil {
			return nil, errors.New("pd is nil")
		}
		tf, err := st.BinarySearch(pd.V2Pair)
		if err != nil {
			return nil, err
		}
		tf.CurrentBlockNumber = artemis_eth_units.NewBigIntFromUint(bn)
		tf.Tx = tx.Tx
		err = ApplyMaxTransferTax(ctx, &tf)
		if err != nil {
			return nil, err
		}
		pend := len(st.Path) - 1
		if m != nil {
			m.StageProgressionMetrics.CountPostProcessTx(float64(1))
			m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactTokensForTokens)
			m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
			m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactTokensForTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount.String(), tf.SandwichPrediction.ExpectedProfit.String())
		}
		log.Info().Msg("swapExactTokensForTokens: saving mempool tx")
		err = SaveMempoolTx(ctx, []web3_client.TradeExecutionFlow{tf}, m)
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
		tfSlice = append(tfSlice, tf)
	case swapTokensForExactTokens:
		st := web3_client.SwapTokensForExactTokensParams{}
		err := st.Decode(tx.Args)
		if err != nil {
			return nil, err
		}
		pend := len(st.Path) - 1
		pd, err := uniswap_pricing.GetV2PricingData(ctx, w3a, st.Path)
		if err != nil {
			if pd != nil && m != nil {
				m.ErrTrackingMetrics.RecordError(swapTokensForExactTokens, pd.V2Pair.PairContractAddr)
			}
			return nil, err
		}
		if pd == nil {
			return nil, errors.New("pd is nil")
		}
		tf, err := st.BinarySearch(pd.V2Pair)
		if err != nil {
			return nil, err
		}
		tf.CurrentBlockNumber = artemis_eth_units.NewBigIntFromUint(bn)
		tf.Tx = tx.Tx
		err = ApplyMaxTransferTax(ctx, &tf)
		if err != nil {
			return nil, err
		}
		if m != nil {
			m.StageProgressionMetrics.CountPostProcessTx(float64(1))

			m.TxFetcherMetrics.TransactionGroup(toAddr, swapTokensForExactTokens)
			m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
			m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapTokensForExactTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount.String(), tf.SandwichPrediction.ExpectedProfit.String())
		}
		log.Info().Msg("swapTokensForExactTokens: saving mempool tx")
		err = SaveMempoolTx(ctx, []web3_client.TradeExecutionFlow{tf}, m)
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
		tfSlice = append(tfSlice, tf)
	case swapExactETHForTokens:
		// payable
		if tx.Tx.Value() == nil {
			return nil, errors.New("swapExactETHForTokens tx has no value")
		}
		st := web3_client.SwapExactETHForTokensParams{}
		err := st.Decode(tx.Args, tx.Tx.Value())
		if err != nil {
			return nil, err
		}
		pend := len(st.Path) - 1
		pd, err := uniswap_pricing.GetV2PricingData(ctx, w3a, st.Path)
		if err != nil {
			if pd != nil && m != nil {
				m.ErrTrackingMetrics.RecordError(swapExactETHForTokens, pd.V2Pair.PairContractAddr)
			}
			return nil, err
		}
		if pd == nil {
			return nil, errors.New("pd is nil")
		}
		tf, err := st.BinarySearch(pd.V2Pair)
		if err != nil {
			return nil, err
		}
		tf.CurrentBlockNumber = artemis_eth_units.NewBigIntFromUint(bn)
		tf.Tx = tx.Tx
		err = ApplyMaxTransferTax(ctx, &tf)
		if err != nil {
			return nil, err
		}
		if m != nil {
			m.StageProgressionMetrics.CountPostProcessTx(float64(1))
			m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactETHForTokens)
			m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
			m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactETHForTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount.String(), tf.SandwichPrediction.ExpectedProfit.String())
		}
		log.Info().Msg("swapExactETHForTokens: saving mempool tx")
		err = SaveMempoolTx(ctx, []web3_client.TradeExecutionFlow{tf}, m)
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
		tfSlice = append(tfSlice, tf)
	case swapTokensForExactETH:
		st := web3_client.SwapTokensForExactETHParams{}
		err := st.Decode(tx.Args)
		if err != nil {
			return nil, err
		}
		pend := len(st.Path) - 1
		pd, err := uniswap_pricing.GetV2PricingData(ctx, w3a, st.Path)
		if err != nil {
			if pd != nil && m != nil {
				m.ErrTrackingMetrics.RecordError(swapTokensForExactETH, pd.V2Pair.PairContractAddr)
			}
			return nil, err
		}
		if pd == nil {
			return nil, errors.New("pd is nil")
		}
		tf, err := st.BinarySearch(pd.V2Pair)
		if err != nil {
			return nil, err
		}
		tf.CurrentBlockNumber = artemis_eth_units.NewBigIntFromUint(bn)
		tf.Tx = tx.Tx
		err = ApplyMaxTransferTax(ctx, &tf)
		if err != nil {
			return nil, err
		}
		log.Info().Msg("swapTokensForExactETH: saving mempool tx")
		err = SaveMempoolTx(ctx, []web3_client.TradeExecutionFlow{tf}, m)
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
		if m != nil {
			m.StageProgressionMetrics.CountPostProcessTx(float64(1))
			m.TxFetcherMetrics.TransactionGroup(toAddr, swapTokensForExactETH)
			m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
			m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapTokensForExactETH, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount.String(), tf.SandwichPrediction.ExpectedProfit.String())
		}
		tfSlice = append(tfSlice, tf)
	case swapExactTokensForETH:
		st := web3_client.SwapExactTokensForETHParams{}
		err := st.Decode(tx.Args)
		if err != nil {
			return nil, err
		}
		pend := len(st.Path) - 1
		pd, err := uniswap_pricing.GetV2PricingData(ctx, w3a, st.Path)
		if err != nil {
			if pd != nil && m != nil {
				m.ErrTrackingMetrics.RecordError(swapExactTokensForETH, pd.V2Pair.PairContractAddr)
			}
			return nil, err
		}
		if pd == nil {
			return nil, errors.New("pd is nil")
		}
		tf, err := st.BinarySearch(pd.V2Pair)
		if err != nil {
			return nil, err
		}
		tf.CurrentBlockNumber = artemis_eth_units.NewBigIntFromUint(bn)
		tf.Tx = tx.Tx
		err = ApplyMaxTransferTax(ctx, &tf)
		if err != nil {
			return nil, err
		}
		log.Info().Msg("swapExactTokensForETH: saving mempool tx")
		err = SaveMempoolTx(ctx, []web3_client.TradeExecutionFlow{tf}, m)
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
		if m != nil {
			m.StageProgressionMetrics.CountPostProcessTx(float64(1))
			m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactTokensForETH)
			m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
			m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactTokensForETH, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount.String(), tf.SandwichPrediction.ExpectedProfit.String())
		}
		tfSlice = append(tfSlice, tf)
	case swapETHForExactTokens:
		// payable
		if tx.Tx.Value() == nil {
			return nil, errors.New("swapETHForExactTokens tx has no value")
		}
		st := web3_client.SwapETHForExactTokensParams{}
		err := st.Decode(tx.Args, tx.Tx.Value())
		if err != nil {
			return nil, err
		}
		pend := len(st.Path) - 1
		pd, err := uniswap_pricing.GetV2PricingData(ctx, w3a, st.Path)
		if err != nil {
			if pd != nil && m != nil {
				m.ErrTrackingMetrics.RecordError(swapETHForExactTokens, pd.V2Pair.PairContractAddr)
			}
			return nil, err
		}
		if pd == nil {
			return nil, errors.New("pd is nil")
		}
		tf, err := st.BinarySearch(pd.V2Pair)
		if err != nil {
			return nil, err
		}
		tf.CurrentBlockNumber = artemis_eth_units.NewBigIntFromUint(bn)
		tf.Tx = tx.Tx
		err = ApplyMaxTransferTax(ctx, &tf)
		if err != nil {
			return nil, err
		}
		log.Info().Msg("swapETHForExactTokens: saving mempool tx")
		err = SaveMempoolTx(ctx, []web3_client.TradeExecutionFlow{tf}, m)
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
		if m != nil {
			m.StageProgressionMetrics.CountPostProcessTx(float64(1))
			m.TxFetcherMetrics.TransactionGroup(toAddr, swapETHForExactTokens)
			m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
			m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapETHForExactTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount.String(), tf.SandwichPrediction.ExpectedProfit.String())
		}
		tfSlice = append(tfSlice, tf)
	}

	if tx.Tx.To().String() != accounts.HexToAddress(web3_client.UniswapV2Router02Address).String() {
		return nil, nil
	}
	switch tx.MethodName {
	case removeLiquidityETHWithPermitSupportingFeeOnTransferTokens:
		if m != nil {
			m.TxFetcherMetrics.TransactionGroup(toAddr, removeLiquidityETHWithPermitSupportingFeeOnTransferTokens)
		}
	case removeLiquidityETHSupportingFeeOnTransferTokens:
		if m != nil {
			m.TxFetcherMetrics.TransactionGroup(toAddr, removeLiquidityETHSupportingFeeOnTransferTokens)
		}
	case swapExactTokensForETHSupportingFeeOnTransferTokens:
		st := web3_client.SwapExactTokensForETHSupportingFeeOnTransferTokensParams{}
		err := st.Decode(tx.Args)
		if err != nil {
			return nil, err
		}
		pend := len(st.Path) - 1
		pd, err := uniswap_pricing.GetV2PricingData(ctx, w3a, st.Path)
		if err != nil {
			if pd != nil && m != nil {
				m.ErrTrackingMetrics.RecordError(swapExactTokensForETHSupportingFeeOnTransferTokens, pd.V2Pair.PairContractAddr)
			}
			return nil, err
		}
		if pd == nil {
			return nil, errors.New("pd is nil")
		}
		tf, err := st.BinarySearch(pd.V2Pair)
		if err != nil {
			return nil, err
		}
		tf.CurrentBlockNumber = artemis_eth_units.NewBigIntFromUint(bn)
		tf.Tx = tx.Tx
		err = ApplyMaxTransferTax(ctx, &tf)
		if err != nil {
			return nil, err
		}
		log.Info().Msg("swapExactTokensForETHSupportingFeeOnTransferTokens: saving mempool tx")
		err = SaveMempoolTx(ctx, []web3_client.TradeExecutionFlow{tf}, m)
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
		if m != nil {
			m.StageProgressionMetrics.CountPostProcessTx(float64(1))
			m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactTokensForETHSupportingFeeOnTransferTokens)
			m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
			m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactTokensForETHSupportingFeeOnTransferTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount.String(), tf.SandwichPrediction.ExpectedProfit.String())
		}
		tfSlice = append(tfSlice, tf)
	case swapExactETHForTokensSupportingFeeOnTransferTokens:
		// payable
		if tx.Tx.Value() == nil {
			return nil, errors.New("swapExactETHForTokensSupportingFeeOnTransferTokens tx has no value")
		}
		st := web3_client.SwapExactETHForTokensSupportingFeeOnTransferTokensParams{}
		err := st.Decode(tx.Args, tx.Tx.Value())
		if err != nil {
			return nil, err
		}
		pend := len(st.Path) - 1
		pd, err := uniswap_pricing.GetV2PricingData(ctx, w3a, st.Path)
		if err != nil {
			if pd != nil && m != nil {
				m.ErrTrackingMetrics.RecordError(swapExactETHForTokensSupportingFeeOnTransferTokens, pd.V2Pair.PairContractAddr)
			}
			return nil, err
		}
		if pd == nil {
			return nil, errors.New("pd is nil")
		}
		tf, err := st.BinarySearch(pd.V2Pair)
		if err != nil {
			return nil, err
		}
		tf.CurrentBlockNumber = artemis_eth_units.NewBigIntFromUint(bn)
		tf.Tx = tx.Tx
		err = ApplyMaxTransferTax(ctx, &tf)
		if err != nil {
			return nil, err
		}
		log.Info().Msg("swapExactETHForTokensSupportingFeeOnTransferTokens: saving mempool tx")
		err = SaveMempoolTx(ctx, []web3_client.TradeExecutionFlow{tf}, m)
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
		if m != nil {
			m.StageProgressionMetrics.CountPostProcessTx(float64(1))
			m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactETHForTokensSupportingFeeOnTransferTokens)
			m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
			m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactETHForTokensSupportingFeeOnTransferTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount.String(), tf.SandwichPrediction.ExpectedProfit.String())
		}
		tfSlice = append(tfSlice, tf)
	case swapExactTokensForTokensSupportingFeeOnTransferTokens:
		st := web3_client.SwapExactTokensForTokensSupportingFeeOnTransferTokensParams{}
		err := st.Decode(tx.Args)
		if err != nil {
			return nil, err
		}
		pend := len(st.Path) - 1
		pd, err := uniswap_pricing.GetV2PricingData(ctx, w3a, st.Path)
		if err != nil {
			if pd != nil && m != nil {
				m.ErrTrackingMetrics.RecordError(swapExactTokensForTokensSupportingFeeOnTransferTokens, pd.V2Pair.PairContractAddr)
			}
			return nil, err
		}
		if pd == nil {
			return nil, errors.New("pd is nil")
		}
		tf, err := st.BinarySearch(pd.V2Pair)
		if err != nil {
			return nil, err
		}
		tf.CurrentBlockNumber = artemis_eth_units.NewBigIntFromUint(bn)
		tf.Tx = tx.Tx
		err = ApplyMaxTransferTax(ctx, &tf)
		if err != nil {
			return nil, err
		}
		log.Info().Msg("swapExactTokensForTokensSupportingFeeOnTransferTokens: saving mempool tx")
		err = SaveMempoolTx(ctx, []web3_client.TradeExecutionFlow{tf}, m)
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
		if m != nil {
			m.StageProgressionMetrics.CountPostProcessTx(float64(1))
			m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactTokensForTokensSupportingFeeOnTransferTokens)
			m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
			m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactTokensForTokensSupportingFeeOnTransferTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount.String(), tf.SandwichPrediction.ExpectedProfit.String())
		}
		tfSlice = append(tfSlice, tf)
	}
	return tfSlice, nil
}
