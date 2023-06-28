package artemis_realtime_trading

import (
	"context"
)

func (a *ActiveTrading) ProcessTxs(ctx context.Context) {
	for _, mevTx := range a.u.MevSmartContractTxMapUniversalRouterOld.Txs {
		a.RealTimeProcessUniversalRouterTx(ctx, mevTx)
	}
	for _, mevTx := range a.u.MevSmartContractTxMapUniversalRouterNew.Txs {
		a.RealTimeProcessUniversalRouterTx(ctx, mevTx)
	}
	//for _, mevTx := range a.u.MevSmartContractTxMapV2Router01.Txs {
	//	a.RealTimeProcessUniswapV2RouterTx(ctx, mevTx)
	//}
	//for _, mevTx := range a.u.MevSmartContractTxMapV2Router02.Txs {
	//	a.RealTimeProcessUniswapV2RouterTx(ctx, mevTx)
	//}
	//for _, mevTx := range a.u.MevSmartContractTxMapV3SwapRouterV2.Txs {
	//	a.RealTimeProcessUniswapV3RouterTx(ctx, mevTx)
	//}
	//for _, mevTx := range a.u.MevSmartContractTxMapV3SwapRouterV1.Txs {
	//	a.RealTimeProcessUniswapV3RouterTx(ctx, mevTx)
	//}
}
