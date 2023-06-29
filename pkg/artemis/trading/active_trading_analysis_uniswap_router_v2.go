package artemis_realtime_trading

import (
	"context"

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
)

func (a *ActiveTrading) RealTimeProcessUniswapV2RouterTx(ctx context.Context, tx web3_client.MevTx) {
	if tx.Tx.To() == nil {
		return
	}
	toAddr := tx.Tx.To().String()
	switch tx.MethodName {
	case addLiquidity:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, addLiquidity)
		//u.AddLiquidity(tx.Args)
	case addLiquidityETH:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, addLiquidityETH)
		// payable
		//u.AddLiquidityETH(tx.Args)
		if tx.Tx.Value() == nil {
			return
		}
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
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactTokensForTokens)
		st := web3_client.SwapExactTokensForTokensParams{}
		st.Decode(ctx, tx.Args)
		pend := len(st.Path) - 1
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
		pd, err := a.u.GetV2PricingData(ctx, st.Path)
		if err != nil {
			return
		}
		tf := st.BinarySearch(pd.V2Pair)
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactTokensForTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
	case swapTokensForExactTokens:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapTokensForExactTokens)
		st := web3_client.SwapTokensForExactTokensParams{}
		st.Decode(tx.Args)
		pend := len(st.Path) - 1
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
		pd, err := a.u.GetV2PricingData(ctx, st.Path)
		if err != nil {
			return
		}
		tf := st.BinarySearch(pd.V2Pair)
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapTokensForExactTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
	case swapExactETHForTokens:
		// payable
		if tx.Tx.Value() == nil {
			return
		}
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactETHForTokens)
		st := web3_client.SwapExactETHForTokensParams{}
		st.Decode(tx.Args, tx.Tx.Value())
		pend := len(st.Path) - 1
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
		pd, err := a.u.GetV2PricingData(ctx, st.Path)
		if err != nil {
			return
		}
		tf := st.BinarySearch(pd.V2Pair)
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactETHForTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
	case swapTokensForExactETH:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapTokensForExactETH)
		st := web3_client.SwapTokensForExactETHParams{}
		st.Decode(tx.Args)
		pend := len(st.Path) - 1
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
		pd, err := a.u.GetV2PricingData(ctx, st.Path)
		if err != nil {
			return
		}
		tf := st.BinarySearch(pd.V2Pair)
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapTokensForExactETH, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
	case swapExactTokensForETH:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactTokensForETH)
		st := web3_client.SwapExactTokensForETHParams{}
		st.Decode(tx.Args)
		pend := len(st.Path) - 1
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
		pd, err := a.u.GetV2PricingData(ctx, st.Path)
		if err != nil {
			return
		}
		tf := st.BinarySearch(pd.V2Pair)
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapExactTokensForETH, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
	case swapETHForExactTokens:
		// payable
		if tx.Tx.Value() == nil {
			return
		}
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapETHForExactTokens)
		st := web3_client.SwapETHForExactTokensParams{}
		st.Decode(tx.Args, tx.Tx.Value())
		pend := len(st.Path) - 1
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
		pd, err := a.u.GetV2PricingData(ctx, st.Path)
		if err != nil {
			return
		}
		tf := st.BinarySearch(pd.V2Pair)
		a.m.TradeAnalysisMetrics.CalculatedSandwichWithPriceLookup(ctx, swapETHForExactTokens, pd.V2Pair.PairContractAddr, st.Path[0].String(), tf.SandwichPrediction.SellAmount, tf.SandwichPrediction.ExpectedProfit)
	}
}
