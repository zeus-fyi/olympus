package artemis_realtime_trading

import (
	"context"
)

func (a *ActiveTrading) ProcessTxs(ctx context.Context) error {
	for _, mevTx := range a.u.MevSmartContractTxMapUniversalRouterOld.Txs {
		err := a.RealTimeProcessUniversalRouterTx(ctx, mevTx)
		if err != nil {
			return err
		}
	}
	for _, mevTx := range a.u.MevSmartContractTxMapUniversalRouterNew.Txs {
		err := a.RealTimeProcessUniversalRouterTx(ctx, mevTx)
		if err != nil {
			return err
		}
	}
	for _, mevTx := range a.u.MevSmartContractTxMapV2Router01.Txs {
		err := a.RealTimeProcessUniswapV2RouterTx(ctx, mevTx)
		if err != nil {
			return err
		}
	}
	for _, mevTx := range a.u.MevSmartContractTxMapV2Router02.Txs {
		err := a.RealTimeProcessUniswapV2RouterTx(ctx, mevTx)
		if err != nil {
			return err
		}
	}
	for _, mevTx := range a.u.MevSmartContractTxMapV3SwapRouterV2.Txs {
		err := a.RealTimeProcessUniswapV3RouterTx(ctx, mevTx, a.u.MevSmartContractTxMapV3SwapRouterV2.Abi, a.u.MevSmartContractTxMapV3SwapRouterV2.Filter)
		if err != nil {
			return err
		}
	}
	for _, mevTx := range a.u.MevSmartContractTxMapV3SwapRouterV1.Txs {
		err := a.RealTimeProcessUniswapV3RouterTx(ctx, mevTx, a.u.MevSmartContractTxMapV3SwapRouterV1.Abi, a.u.MevSmartContractTxMapV3SwapRouterV1.Filter)
		if err != nil {
			return err
		}
	}
	return nil
}
