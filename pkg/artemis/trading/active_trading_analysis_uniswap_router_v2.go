package artemis_realtime_trading

import (
	"context"
	"errors"

	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

const (
	addLiquidity                                       = "addLiquidity"
	addLiquidityETH                                    = "addLiquidityETH"
	removeLiquidity                                    = "removeLiquidity"
	removeLiquidityETH                                 = "removeLiquidityETH"
	removeLiquidityWithPermit                          = "removeLiquidityWithPermit"
	removeLiquidityETHWithPermit                       = "removeLiquidityETHWithPermit"
	swapExactTokensForTokens                           = "swapExactTokensForTokens"
	swapTokensForExactTokens                           = "swapTokensForExactTokens"
	swapExactETHForTokens                              = "swapExactETHForTokens"
	swapTokensForExactETH                              = "swapTokensForExactETH"
	swapExactTokensForETH                              = "swapExactTokensForETH"
	swapETHForExactTokens                              = "swapETHForExactTokens"
	swapExactTokensForETHSupportingFeeOnTransferTokens = "swapExactTokensForETHSupportingFeeOnTransferTokens"
	swapExactETHForTokensSupportingFeeOnTransferTokens = "swapExactETHForTokensSupportingFeeOnTransferTokens"
)

func (a *ActiveTrading) RealTimeProcessUniswapV2RouterTx(ctx context.Context, tx web3_client.MevTx) error {
	toAddr := tx.Tx.To().String()
	switch tx.MethodName {
	case addLiquidity:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, addLiquidity)
		//u.AddLiquidity(tx.Args)
	case addLiquidityETH:
		// payable
		//u.AddLiquidityETH(tx.Args)
		if tx.Tx.Value() == nil {
			return errors.New("addLiquidityETH tx has no value")
		}
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, addLiquidityETH)
	case removeLiquidity:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, removeLiquidity)
		//u.RemoveLiquidity(tx.Args)
	case removeLiquidityETH:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, removeLiquidityETH)
		//u.RemoveLiquidityETH(tx.Args)
	case removeLiquidityWithPermit:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, removeLiquidityWithPermit)
		//u.RemoveLiquidityWithPermit(tx.Args)
	case removeLiquidityETHWithPermit:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, removeLiquidityETHWithPermit)
		//u.RemoveLiquidityETHWithPermit(tx.Args)
	case swapExactTokensForTokens:
		st := web3_client.SwapExactTokensForTokensParams{}
		st.Decode(ctx, tx.Args)
		pd, err := a.u.GetV2PricingData(ctx, st.Path)
		if err != nil {
			return err
		}
		tf := st.BinarySearch(pd.V2Pair)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" {
			return errors.New("expectedProfit == 0 or 1")
		}
		pend := len(st.Path) - 1
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactTokensForTokens)
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactTokensForTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
	case swapTokensForExactTokens:
		st := web3_client.SwapTokensForExactTokensParams{}
		st.Decode(tx.Args)
		pend := len(st.Path) - 1
		pd, err := a.u.GetV2PricingData(ctx, st.Path)
		if err != nil {
			return err
		}
		tf := st.BinarySearch(pd.V2Pair)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" {
			return errors.New("expectedProfit == 0 or 1")
		}
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapTokensForExactTokens)
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapTokensForExactTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
	case swapExactETHForTokens:
		// payable
		if tx.Tx.Value() == nil {
			return errors.New("swapExactETHForTokens tx has no value")
		}
		st := web3_client.SwapExactETHForTokensParams{}
		st.Decode(tx.Args, tx.Tx.Value())
		pend := len(st.Path) - 1
		pd, err := a.u.GetV2PricingData(ctx, st.Path)
		if err != nil {
			return err
		}
		tf := st.BinarySearch(pd.V2Pair)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" {
			return errors.New("expectedProfit == 0 or 1")
		}
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactETHForTokens)
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactETHForTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
	case swapTokensForExactETH:
		st := web3_client.SwapTokensForExactETHParams{}
		st.Decode(tx.Args)
		pend := len(st.Path) - 1
		pd, err := a.u.GetV2PricingData(ctx, st.Path)
		if err != nil {
			return err
		}
		tf := st.BinarySearch(pd.V2Pair)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" {
			return errors.New("expectedProfit == 0 or 1")
		}
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapTokensForExactETH)
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapTokensForExactETH, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
	case swapExactTokensForETH:
		st := web3_client.SwapExactTokensForETHParams{}
		st.Decode(tx.Args)
		pend := len(st.Path) - 1
		pd, err := a.u.GetV2PricingData(ctx, st.Path)
		if err != nil {
			return err
		}
		tf := st.BinarySearch(pd.V2Pair)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" {
			return errors.New("expectedProfit == 0 or 1")
		}
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactTokensForETH)
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactTokensForETH, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
	case swapETHForExactTokens:
		// payable
		if tx.Tx.Value() == nil {
			return errors.New("swapETHForExactTokens tx has no value")
		}
		st := web3_client.SwapETHForExactTokensParams{}
		st.Decode(tx.Args, tx.Tx.Value())
		pend := len(st.Path) - 1
		pd, err := a.u.GetV2PricingData(ctx, st.Path)
		if err != nil {
			return err
		}
		tf := st.BinarySearch(pd.V2Pair)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" {
			return errors.New("expectedProfit == 0 or 1")
		}
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapETHForExactTokens)
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapETHForExactTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
	case swapExactTokensForETHSupportingFeeOnTransferTokens:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactTokensForETHSupportingFeeOnTransferTokens)
		st := web3_client.SwapExactTokensForETHSupportingFeeOnTransferTokensParams{}
		st.Decode(tx.Args)
		pend := len(st.Path) - 1
		pd, err := a.u.GetV2PricingData(ctx, st.Path)
		if err != nil {
			return err
		}
		tf := st.BinarySearch(pd.V2Pair)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" {
			return errors.New("expectedProfit == 0 or 1")
		}
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactTokensForETHSupportingFeeOnTransferTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
	case swapExactETHForTokensSupportingFeeOnTransferTokens:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactETHForTokensSupportingFeeOnTransferTokens)
		// payable
		if tx.Tx.Value() == nil {
			return errors.New("swapExactETHForTokensSupportingFeeOnTransferTokens tx has no value")
		}
		st := web3_client.SwapExactETHForTokensSupportingFeeOnTransferTokensParams{}
		st.Decode(tx.Args, tx.Tx.Value())
		pend := len(st.Path) - 1
		pd, err := a.u.GetV2PricingData(ctx, st.Path)
		if err != nil {
			return err
		}
		tf := st.BinarySearch(pd.V2Pair)
		if tf.SandwichPrediction.ExpectedProfit == "0" || tf.SandwichPrediction.ExpectedProfit == "1" {
			return errors.New("expectedProfit == 0 or 1")
		}
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactETHForTokensSupportingFeeOnTransferTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
	}
	return nil
}
