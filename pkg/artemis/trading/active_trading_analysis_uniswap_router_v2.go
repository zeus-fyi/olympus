package artemis_realtime_trading

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
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

func (a *ActiveTrading) RealTimeProcessUniswapV2RouterTx(ctx context.Context, tx web3_client.MevTx) ([]web3_client.TradeExecutionFlowJSON, error) {
	toAddr := tx.Tx.To().String()
	var tfSlice []web3_client.TradeExecutionFlowJSON
	switch tx.MethodName {
	case addLiquidity:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, addLiquidity)
	case addLiquidityETH:
		if tx.Tx.Value() == nil {
			return nil, errors.New("addLiquidityETH tx has no value")
		}
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, addLiquidityETH)
	case removeLiquidity:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, removeLiquidity)
	case removeLiquidityETH:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, removeLiquidityETH)
	case removeLiquidityWithPermit:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, removeLiquidityWithPermit)
	case removeLiquidityETHWithPermit:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, removeLiquidityETHWithPermit)
	case swapExactTokensForTokens:
		st := web3_client.SwapExactTokensForTokensParams{}
		st.Decode(ctx, tx.Args)
		pd, err := a.u.GetV2PricingData(ctx, st.Path)
		if err != nil {
			a.m.ErrTrackingMetrics.RecordError(swapExactTokensForTokens, pd.V2Pair.PairContractAddr)
			return nil, err
		}
		tf := st.BinarySearch(pd.V2Pair)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" {
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
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactTokensForTokens)
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactTokensForTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		tfSlice = append(tfSlice, tf)
	case swapTokensForExactTokens:
		st := web3_client.SwapTokensForExactTokensParams{}
		st.Decode(tx.Args)
		pend := len(st.Path) - 1
		pd, err := a.u.GetV2PricingData(ctx, st.Path)
		if err != nil {
			a.m.ErrTrackingMetrics.RecordError(swapTokensForExactTokens, pd.V2Pair.PairContractAddr)
			return nil, err
		}
		tf := st.BinarySearch(pd.V2Pair)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" {
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
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapTokensForExactTokens)
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapTokensForExactTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		tfSlice = append(tfSlice, tf)
	case swapExactETHForTokens:
		// payable
		if tx.Tx.Value() == nil {
			return nil, errors.New("swapExactETHForTokens tx has no value")
		}
		st := web3_client.SwapExactETHForTokensParams{}
		st.Decode(tx.Args, tx.Tx.Value())
		pend := len(st.Path) - 1
		pd, err := a.u.GetV2PricingData(ctx, st.Path)
		if err != nil {
			a.m.ErrTrackingMetrics.RecordError(swapExactETHForTokens, pd.V2Pair.PairContractAddr)
			return nil, err
		}
		tf := st.BinarySearch(pd.V2Pair)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" {
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
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactETHForTokens)
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactETHForTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		tfSlice = append(tfSlice, tf)
	case swapTokensForExactETH:
		st := web3_client.SwapTokensForExactETHParams{}
		st.Decode(tx.Args)
		pend := len(st.Path) - 1
		pd, err := a.u.GetV2PricingData(ctx, st.Path)
		if err != nil {
			a.m.ErrTrackingMetrics.RecordError(swapTokensForExactETH, pd.V2Pair.PairContractAddr)
			return nil, err
		}
		tf := st.BinarySearch(pd.V2Pair)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" {
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
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapTokensForExactETH)
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapTokensForExactETH, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		tfSlice = append(tfSlice, tf)
	case swapExactTokensForETH:
		st := web3_client.SwapExactTokensForETHParams{}
		st.Decode(tx.Args)
		pend := len(st.Path) - 1
		pd, err := a.u.GetV2PricingData(ctx, st.Path)
		if err != nil {
			a.m.ErrTrackingMetrics.RecordError(swapExactTokensForETH, pd.V2Pair.PairContractAddr)
			return nil, err
		}
		tf := st.BinarySearch(pd.V2Pair)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" {
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
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactTokensForETH)
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactTokensForETH, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		tfSlice = append(tfSlice, tf)
	case swapETHForExactTokens:
		// payable
		if tx.Tx.Value() == nil {
			return nil, errors.New("swapETHForExactTokens tx has no value")
		}
		st := web3_client.SwapETHForExactTokensParams{}
		st.Decode(tx.Args, tx.Tx.Value())
		pend := len(st.Path) - 1
		pd, err := a.u.GetV2PricingData(ctx, st.Path)
		if err != nil {
			a.m.ErrTrackingMetrics.RecordError(swapETHForExactTokens, pd.V2Pair.PairContractAddr)
			return nil, err
		}
		tf := st.BinarySearch(pd.V2Pair)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" {
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
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapETHForExactTokens)
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapETHForExactTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		tfSlice = append(tfSlice, tf)
	}

	if tx.Tx.To().String() != accounts.HexToAddress(web3_client.UniswapV2Router02Address).String() {
		return nil, nil
	}
	switch tx.MethodName {
	case removeLiquidityETHWithPermitSupportingFeeOnTransferTokens:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, removeLiquidityETHWithPermitSupportingFeeOnTransferTokens)
	case removeLiquidityETHSupportingFeeOnTransferTokens:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, removeLiquidityETHSupportingFeeOnTransferTokens)
	case swapExactTokensForETHSupportingFeeOnTransferTokens:
		st := web3_client.SwapExactTokensForETHSupportingFeeOnTransferTokensParams{}
		st.Decode(tx.Args)
		pend := len(st.Path) - 1
		pd, err := a.u.GetV2PricingData(ctx, st.Path)
		if err != nil {
			a.m.ErrTrackingMetrics.RecordError(swapExactTokensForETHSupportingFeeOnTransferTokens, pd.V2Pair.PairContractAddr)
			return nil, err
		}
		tf := st.BinarySearch(pd.V2Pair)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" {
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
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactTokensForETHSupportingFeeOnTransferTokens)
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactTokensForETHSupportingFeeOnTransferTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		tfSlice = append(tfSlice, tf)
	case swapExactETHForTokensSupportingFeeOnTransferTokens:
		// payable
		if tx.Tx.Value() == nil {
			return nil, errors.New("swapExactETHForTokensSupportingFeeOnTransferTokens tx has no value")
		}
		st := web3_client.SwapExactETHForTokensSupportingFeeOnTransferTokensParams{}
		st.Decode(tx.Args, tx.Tx.Value())
		pend := len(st.Path) - 1
		pd, err := a.u.GetV2PricingData(ctx, st.Path)
		if err != nil {
			a.m.ErrTrackingMetrics.RecordError(swapExactETHForTokensSupportingFeeOnTransferTokens, pd.V2Pair.PairContractAddr)
			return nil, err
		}
		tf := st.BinarySearch(pd.V2Pair)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" {
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
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactETHForTokensSupportingFeeOnTransferTokens)
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactETHForTokensSupportingFeeOnTransferTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		tfSlice = append(tfSlice, tf)
	case swapExactTokensForTokensSupportingFeeOnTransferTokens:
		st := web3_client.SwapExactTokensForTokensSupportingFeeOnTransferTokensParams{}
		st.Decode(tx.Args)
		pend := len(st.Path) - 1
		pd, err := a.u.GetV2PricingData(ctx, st.Path)
		if err != nil {
			a.m.ErrTrackingMetrics.RecordError(swapExactTokensForTokensSupportingFeeOnTransferTokens, pd.V2Pair.PairContractAddr)
			return nil, err
		}
		tf := st.BinarySearch(pd.V2Pair)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" {
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
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactTokensForTokensSupportingFeeOnTransferTokens)
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactTokensForTokensSupportingFeeOnTransferTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
		tfSlice = append(tfSlice, tf)
	}
	return tfSlice, nil
}
