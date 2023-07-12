package artemis_realtime_trading

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *ActiveTrading) DecodeTx(ctx context.Context, tx *types.Transaction) ([]web3_client.MevTx, error) {
	var mevTxs []web3_client.MevTx
	switch tx.To().String() {
	case web3_client.UniswapUniversalRouterAddressOld:
		if a.a.U.MevSmartContractTxMapUniversalRouterOld.Txs == nil {
			a.a.U.MevSmartContractTxMapUniversalRouterOld.Txs = []web3_client.MevTx{}
		}
		methodName, args, err := web3_client.DecodeTxArgData(ctx, tx, a.a.U.MevSmartContractTxMapUniversalRouterOld)
		if err != nil {
			return nil, err
		}
		singleTx := web3_client.MevTx{
			MethodName: methodName,
			//UserAddr:    userAddr,
			Args: args,
			//Order:       order,
			//TxPoolQueue: txPoolQueue,
			Tx: tx,
		}
		tmp := a.a.U.MevSmartContractTxMapUniversalRouterOld.Txs
		tmp = append(tmp, singleTx)
		mevTxs = append(mevTxs, singleTx)
		a.a.U.MevSmartContractTxMapUniversalRouterOld.Txs = tmp
		a.m.TxFetcherMetrics.TransactionGroup(web3_client.UniswapUniversalRouterAddressOld, methodName)
	case web3_client.UniswapUniversalRouterAddressNew:
		if a.a.U.MevSmartContractTxMapUniversalRouterNew.Txs == nil {
			a.a.U.MevSmartContractTxMapUniversalRouterNew.Txs = []web3_client.MevTx{}
		}
		methodName, args, err := web3_client.DecodeTxArgData(ctx, tx, a.a.U.MevSmartContractTxMapUniversalRouterNew)
		if err != nil {
			return nil, err
		}
		singleTx := web3_client.MevTx{
			MethodName: methodName,
			//UserAddr:    userAddr,
			Args: args,
			//Order:       order,
			//TxPoolQueue: txPoolQueue,
			Tx: tx,
		}
		tmp := a.a.U.MevSmartContractTxMapUniversalRouterNew.Txs
		tmp = append(tmp, singleTx)
		mevTxs = append(mevTxs, singleTx)
		a.a.U.MevSmartContractTxMapUniversalRouterNew.Txs = tmp
		a.m.TxFetcherMetrics.TransactionGroup(web3_client.UniswapUniversalRouterAddressNew, methodName)
	case web3_client.UniswapV2Router02Address:
		if a.a.U.MevSmartContractTxMapV2Router02.Txs == nil {
			a.a.U.MevSmartContractTxMapV2Router02.Txs = []web3_client.MevTx{}
		}
		methodName, args, err := web3_client.DecodeTxArgData(ctx, tx, a.a.U.MevSmartContractTxMapV2Router02)
		if err != nil {
			return nil, err
		}
		singleTx := web3_client.MevTx{
			MethodName: methodName,
			//UserAddr:    userAddr,
			Args: args,
			//Order:       order,
			//TxPoolQueue: txPoolQueue,
			Tx: tx,
		}
		tmp := a.a.U.MevSmartContractTxMapV2Router02.Txs
		tmp = append(tmp, singleTx)
		mevTxs = append(mevTxs, singleTx)
		a.a.U.MevSmartContractTxMapV2Router02.Txs = tmp
		a.m.TxFetcherMetrics.TransactionGroup(web3_client.UniswapV2Router02Address, methodName)
	case web3_client.UniswapV2Router01Address:
		if a.a.U.MevSmartContractTxMapV2Router01.Txs == nil {
			a.a.U.MevSmartContractTxMapV2Router01.Txs = []web3_client.MevTx{}
		}
		methodName, args, err := web3_client.DecodeTxArgData(ctx, tx, a.a.U.MevSmartContractTxMapV2Router01)
		if err != nil {
			return nil, err
		}
		singleTx := web3_client.MevTx{
			MethodName: methodName,
			//UserAddr:    userAddr,
			Args: args,
			//Order:       order,
			//TxPoolQueue: txPoolQueue,
			Tx: tx,
		}
		tmp := a.a.U.MevSmartContractTxMapV2Router01.Txs
		tmp = append(tmp, singleTx)
		mevTxs = append(mevTxs, singleTx)
		a.a.U.MevSmartContractTxMapV2Router01.Txs = tmp
		a.m.TxFetcherMetrics.TransactionGroup(web3_client.UniswapV2Router01Address, methodName)
	case web3_client.UniswapV3Router01Address:
		if a.a.U.MevSmartContractTxMapV3SwapRouterV1.Txs == nil {
			a.a.U.MevSmartContractTxMapV3SwapRouterV1.Txs = []web3_client.MevTx{}
		}
		methodName, args, err := web3_client.DecodeTxArgData(ctx, tx, a.a.U.MevSmartContractTxMapV3SwapRouterV1)
		if err != nil {
			return nil, err
		}
		singleTx := web3_client.MevTx{
			MethodName: methodName,
			//UserAddr:    userAddr,
			Args: args,
			//Order:       order,
			//TxPoolQueue: txPoolQueue,
			Tx: tx,
		}
		tmp := a.a.U.MevSmartContractTxMapV3SwapRouterV1.Txs
		tmp = append(tmp, singleTx)
		mevTxs = append(mevTxs, singleTx)
		a.a.U.MevSmartContractTxMapV3SwapRouterV1.Txs = tmp
		a.m.TxFetcherMetrics.TransactionGroup(web3_client.UniswapV3Router01Address, methodName)
	case web3_client.UniswapV3Router02Address:
		if a.a.U.MevSmartContractTxMapV3SwapRouterV2.Txs == nil {
			a.a.U.MevSmartContractTxMapV3SwapRouterV2.Txs = []web3_client.MevTx{}
		}
		methodName, args, err := web3_client.DecodeTxArgData(ctx, tx, a.a.U.MevSmartContractTxMapV3SwapRouterV2)
		if err != nil {
			return nil, err
		}
		singleTx := web3_client.MevTx{
			MethodName: methodName,
			//UserAddr:    userAddr,
			Args: args,
			//Order:       order,
			//TxPoolQueue: txPoolQueue,
			Tx: tx,
		}
		tmp := a.a.U.MevSmartContractTxMapV3SwapRouterV2.Txs
		tmp = append(tmp, singleTx)
		mevTxs = append(mevTxs, singleTx)
		a.a.U.MevSmartContractTxMapV3SwapRouterV2.Txs = tmp
		a.m.TxFetcherMetrics.TransactionGroup(web3_client.UniswapV3Router02Address, methodName)
	}
	return mevTxs, nil
}
