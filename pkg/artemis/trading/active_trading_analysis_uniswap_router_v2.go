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
		//u.AddLiquidity(tx.Args)
	case addLiquidityETH:
		// payable
		//u.AddLiquidityETH(tx.Args)
		if tx.Tx.Value() == nil {
			return
		}
	case removeLiquidity:
		//u.RemoveLiquidity(tx.Args)
	case removeLiquidityETH:
		//u.RemoveLiquidityETH(tx.Args)
	case removeLiquidityWithPermit:
		//u.RemoveLiquidityWithPermit(tx.Args)
	case removeLiquidityETHWithPermit:
		//u.RemoveLiquidityETHWithPermit(tx.Args)
	case swapExactTokensForTokens:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactTokensForTokens)
		st := web3_client.SwapExactTokensForTokensParams{}
		st.Decode(ctx, tx.Args)
		pend := len(st.Path) - 1
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
	case swapTokensForExactTokens:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapTokensForExactTokens)
		st := web3_client.SwapTokensForExactTokensParams{}
		st.Decode(tx.Args)
		pend := len(st.Path) - 1
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
	case swapExactETHForTokens:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactETHForTokens)
		// payable
		if tx.Tx.Value() == nil {
			return
		}
		st := web3_client.SwapExactETHForTokensParams{}
		st.Decode(tx.Args, tx.Tx.Value())
		pend := len(st.Path) - 1
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
	case swapTokensForExactETH:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapTokensForExactETH)
		//a.u.SwapTokensForExactETH(tx, tx.Args)
		st := web3_client.SwapTokensForExactETHParams{}
		st.Decode(tx.Args)
		pend := len(st.Path) - 1
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
	case swapExactTokensForETH:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactTokensForETH)
		//a.u.SwapExactTokensForETH(tx, tx.Args)
		st := web3_client.SwapExactTokensForETHParams{}
		st.Decode(tx.Args)
		pend := len(st.Path) - 1
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
	case swapETHForExactTokens:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapETHForExactTokens)
		// payable
		if tx.Tx.Value() == nil {
			return
		}
		//a.u.SwapETHForExactTokens(tx, tx.Args, tx.Tx.Value())
		st := web3_client.SwapETHForExactTokensParams{}
		st.Decode(tx.Args, tx.Tx.Value())
		pend := len(st.Path) - 1
		a.m.TxFetcherMetrics.TransactionCurrencyInOut(toAddr, st.Path[0].String(), st.Path[pend].String())
	}
}
