package artemis_realtime_trading

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	metrics_trading "github.com/zeus-fyi/olympus/pkg/apollo/ethereum/mev/trading"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
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

func RealTimeProcessUniswapV2RouterTx(ctx context.Context, tx web3_client.MevTx, m *metrics_trading.TradingMetrics, w3a web3_actions.Web3Actions) ([]web3_client.TradeExecutionFlowJSON, error) {
	bn, berr := artemis_trading_cache.GetLatestBlock(ctx)
	if berr != nil {
		log.Err(berr).Msg("failed to get latest block")
		return nil, errors.New("ailed to get latest block")
	}
	toAddr := tx.Tx.To().String()
	var tfSlice []web3_client.TradeExecutionFlowJSON
	switch tx.MethodName {
	case addLiquidity:
		if m != nil {
			m.TxFetcherMetrics.TransactionGroup(toAddr, addLiquidity)
		}
	case addLiquidityETH:
		if tx.Tx.Value() == nil {
			return nil, errors.New("addLiquidityETH tx has no value")
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
		st.Decode(ctx, tx.Args)
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
		tf := st.BinarySearch(pd.V2Pair)
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
		pend := len(st.Path) - 1
		if m != nil {
			m.StageProgressionMetrics.CountPostProcessTx(float64(1))
			m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactTokensForTokens)
			m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
			m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactTokensForTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		}
		tfSlice = append(tfSlice, tf)
		log.Info().Msg("saving mempool tx")
		err = SaveMempoolTx(ctx, bn, []web3_client.TradeExecutionFlowJSON{tf}, m)
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
	case swapTokensForExactTokens:
		st := web3_client.SwapTokensForExactTokensParams{}
		st.Decode(tx.Args)
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
		tf := st.BinarySearch(pd.V2Pair)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" || tf.SandwichPrediction.ExpectedProfit == "" {
			return nil, errors.New("expectedProfit == 0 or 1")
		}
		tf.Trade.TradeMethod = swapTokensForExactTokens
		newTx := artemis_trading_types.JSONTx{}
		err = newTx.UnmarshalTx(tx.Tx)
		if err != nil {
			log.Err(err).Msg("failed to unmarshal tx")
			return nil, err
		}
		tf.Tx = newTx
		tf.InitialPair = pd.V2Pair.ConvertToJSONType()
		if m != nil {
			m.StageProgressionMetrics.CountPostProcessTx(float64(1))

			m.TxFetcherMetrics.TransactionGroup(toAddr, swapTokensForExactTokens)
			m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
			m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapTokensForExactTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		}
		tfSlice = append(tfSlice, tf)
		log.Info().Msg("saving mempool tx")
		err = SaveMempoolTx(ctx, bn, []web3_client.TradeExecutionFlowJSON{tf}, m)
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
	case swapExactETHForTokens:
		// payable
		if tx.Tx.Value() == nil {
			return nil, errors.New("swapExactETHForTokens tx has no value")
		}
		st := web3_client.SwapExactETHForTokensParams{}
		st.Decode(tx.Args, tx.Tx.Value())
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
		tf := st.BinarySearch(pd.V2Pair)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" || tf.SandwichPrediction.ExpectedProfit == "" {
			return nil, errors.New("expectedProfit == 0 or 1")
		}
		tf.Trade.TradeMethod = swapExactETHForTokens
		newTx := artemis_trading_types.JSONTx{}
		err = newTx.UnmarshalTx(tx.Tx)
		if err != nil {
			log.Err(err).Msg("failed to unmarshal tx")
			return nil, err
		}
		tf.Tx = newTx
		tf.InitialPair = pd.V2Pair.ConvertToJSONType()
		if m != nil {
			m.StageProgressionMetrics.CountPostProcessTx(float64(1))
			m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactETHForTokens)
			m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
			m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactETHForTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		}
		tfSlice = append(tfSlice, tf)
		log.Info().Msg("saving mempool tx")
		err = SaveMempoolTx(ctx, bn, []web3_client.TradeExecutionFlowJSON{tf}, m)
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
	case swapTokensForExactETH:
		st := web3_client.SwapTokensForExactETHParams{}
		st.Decode(tx.Args)
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
		tf := st.BinarySearch(pd.V2Pair)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" || tf.SandwichPrediction.ExpectedProfit == "" {
			return nil, errors.New("expectedProfit == 0 or 1")
		}
		tf.Trade.TradeMethod = swapTokensForExactETH
		newTx := artemis_trading_types.JSONTx{}
		err = newTx.UnmarshalTx(tx.Tx)
		if err != nil {
			log.Err(err).Msg("failed to unmarshal tx")
			return nil, err
		}
		tf.Tx = newTx
		tf.InitialPair = pd.V2Pair.ConvertToJSONType()
		if m != nil {
			m.StageProgressionMetrics.CountPostProcessTx(float64(1))
			m.TxFetcherMetrics.TransactionGroup(toAddr, swapTokensForExactETH)
			m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
			m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapTokensForExactETH, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		}
		tfSlice = append(tfSlice, tf)
		log.Info().Msg("saving mempool tx")
		err = SaveMempoolTx(ctx, bn, []web3_client.TradeExecutionFlowJSON{tf}, m)
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
	case swapExactTokensForETH:
		st := web3_client.SwapExactTokensForETHParams{}
		st.Decode(tx.Args)
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
		tf := st.BinarySearch(pd.V2Pair)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" || tf.SandwichPrediction.ExpectedProfit == "" {
			return nil, errors.New("expectedProfit == 0 or 1")
		}
		tf.Trade.TradeMethod = swapExactTokensForETH
		newTx := artemis_trading_types.JSONTx{}
		err = newTx.UnmarshalTx(tx.Tx)
		if err != nil {
			log.Err(err).Msg("failed to unmarshal tx")
			return nil, err
		}
		tf.Tx = newTx
		tf.InitialPair = pd.V2Pair.ConvertToJSONType()
		if m != nil {
			m.StageProgressionMetrics.CountPostProcessTx(float64(1))
			m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactTokensForETH)
			m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
			m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactTokensForETH, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		}
		tfSlice = append(tfSlice, tf)
		log.Info().Msg("saving mempool tx")
		err = SaveMempoolTx(ctx, bn, []web3_client.TradeExecutionFlowJSON{tf}, m)
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
	case swapETHForExactTokens:
		// payable
		if tx.Tx.Value() == nil {
			return nil, errors.New("swapETHForExactTokens tx has no value")
		}
		st := web3_client.SwapETHForExactTokensParams{}
		st.Decode(tx.Args, tx.Tx.Value())
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
		tf := st.BinarySearch(pd.V2Pair)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" || tf.SandwichPrediction.ExpectedProfit == "" {
			return nil, errors.New("expectedProfit == 0 or 1")
		}
		tf.Trade.TradeMethod = swapETHForExactTokens
		newTx := artemis_trading_types.JSONTx{}
		err = newTx.UnmarshalTx(tx.Tx)
		if err != nil {
			log.Err(err).Msg("failed to unmarshal tx")
			return nil, err
		}
		tf.Tx = newTx
		tf.InitialPair = pd.V2Pair.ConvertToJSONType()
		if m != nil {
			m.StageProgressionMetrics.CountPostProcessTx(float64(1))
			m.TxFetcherMetrics.TransactionGroup(toAddr, swapETHForExactTokens)
			m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
			m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapETHForExactTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		}
		tfSlice = append(tfSlice, tf)
		log.Info().Msg("saving mempool tx")
		err = SaveMempoolTx(ctx, bn, []web3_client.TradeExecutionFlowJSON{tf}, m)
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
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
		st.Decode(tx.Args)
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
		tf := st.BinarySearch(pd.V2Pair)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" || tf.SandwichPrediction.ExpectedProfit == "" {
			return nil, errors.New("expectedProfit == 0 or 1")
		}
		tf.Trade.TradeMethod = swapExactTokensForETHSupportingFeeOnTransferTokens
		newTx := artemis_trading_types.JSONTx{}
		err = newTx.UnmarshalTx(tx.Tx)
		if err != nil {
			log.Err(err).Msg("failed to unmarshal tx")
			return nil, err
		}
		tf.Tx = newTx
		tf.InitialPair = pd.V2Pair.ConvertToJSONType()
		if m != nil {
			m.StageProgressionMetrics.CountPostProcessTx(float64(1))
			m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactTokensForETHSupportingFeeOnTransferTokens)
			m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
			m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactTokensForETHSupportingFeeOnTransferTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		}
		tfSlice = append(tfSlice, tf)
		log.Info().Msg("saving mempool tx")
		err = SaveMempoolTx(ctx, bn, []web3_client.TradeExecutionFlowJSON{tf}, m)
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
	case swapExactETHForTokensSupportingFeeOnTransferTokens:
		// payable
		if tx.Tx.Value() == nil {
			return nil, errors.New("swapExactETHForTokensSupportingFeeOnTransferTokens tx has no value")
		}
		st := web3_client.SwapExactETHForTokensSupportingFeeOnTransferTokensParams{}
		st.Decode(tx.Args, tx.Tx.Value())
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
		tf := st.BinarySearch(pd.V2Pair)
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
		tf.Trade.TradeMethod = swapExactETHForTokensSupportingFeeOnTransferTokens
		tf.InitialPair = pd.V2Pair.ConvertToJSONType()
		if m != nil {
			m.StageProgressionMetrics.CountPostProcessTx(float64(1))
			m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactETHForTokensSupportingFeeOnTransferTokens)
			m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
			m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactETHForTokensSupportingFeeOnTransferTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		}
		tfSlice = append(tfSlice, tf)

		log.Info().Msg("saving mempool tx")
		err = SaveMempoolTx(ctx, bn, []web3_client.TradeExecutionFlowJSON{tf}, m)
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
	case swapExactTokensForTokensSupportingFeeOnTransferTokens:
		st := web3_client.SwapExactTokensForTokensSupportingFeeOnTransferTokensParams{}
		st.Decode(tx.Args)
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
		tf := st.BinarySearch(pd.V2Pair)
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
		tf.Trade.TradeMethod = swapExactTokensForTokensSupportingFeeOnTransferTokens
		tf.InitialPair = pd.V2Pair.ConvertToJSONType()
		if m != nil {
			m.StageProgressionMetrics.CountPostProcessTx(float64(1))
			m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactTokensForTokensSupportingFeeOnTransferTokens)
			m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
			m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactTokensForTokensSupportingFeeOnTransferTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		}
		tfSlice = append(tfSlice, tf)
		log.Info().Msg("saving mempool tx")
		err = SaveMempoolTx(ctx, bn, []web3_client.TradeExecutionFlowJSON{tf}, m)
		if err != nil {
			log.Err(err).Msg("failed to save mempool tx")
			return nil, errors.New("failed to save mempool tx")
		}
	}
	return tfSlice, nil
}
