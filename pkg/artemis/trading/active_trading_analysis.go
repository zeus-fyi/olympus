package artemis_realtime_trading

import (
	"context"

	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *ActiveTrading) ProcessTxs(ctx context.Context) ([]*web3_client.TradeExecutionFlowJSON, error) {
	var tfSlice []*web3_client.TradeExecutionFlowJSON
	for _, mevTx := range a.u.MevSmartContractTxMapUniversalRouterOld.Txs {
		tf, err := a.RealTimeProcessUniversalRouterTx(ctx, mevTx)
		if err != nil {
			return nil, err
		}
		tfSlice = append(tfSlice, tf...)
	}
	for _, mevTx := range a.u.MevSmartContractTxMapUniversalRouterNew.Txs {
		tf, err := a.RealTimeProcessUniversalRouterTx(ctx, mevTx)
		if err != nil {
			return nil, err
		}
		tfSlice = append(tfSlice, tf...)
	}
	for _, mevTx := range a.u.MevSmartContractTxMapV2Router01.Txs {
		tf, err := a.RealTimeProcessUniswapV2RouterTx(ctx, mevTx)
		if err != nil {
			return nil, err
		}
		tfSlice = append(tfSlice, tf...)
	}
	for _, mevTx := range a.u.MevSmartContractTxMapV2Router02.Txs {
		tf, err := a.RealTimeProcessUniswapV2RouterTx(ctx, mevTx)
		if err != nil {
			return nil, err
		}
		tfSlice = append(tfSlice, tf...)
	}
	for _, mevTx := range a.u.MevSmartContractTxMapV3SwapRouterV2.Txs {
		tf, err := a.RealTimeProcessUniswapV3RouterTx(ctx, mevTx, a.u.MevSmartContractTxMapV3SwapRouterV2.Abi, a.u.MevSmartContractTxMapV3SwapRouterV2.Filter)
		if err != nil {
			return nil, err
		}
		tfSlice = append(tfSlice, tf...)
	}
	for _, mevTx := range a.u.MevSmartContractTxMapV3SwapRouterV1.Txs {
		tf, err := a.RealTimeProcessUniswapV3RouterTx(ctx, mevTx, a.u.MevSmartContractTxMapV3SwapRouterV1.Abi, a.u.MevSmartContractTxMapV3SwapRouterV1.Filter)
		if err != nil {
			return nil, err
		}
		tfSlice = append(tfSlice, tf...)
	}
	return tfSlice, nil
}
