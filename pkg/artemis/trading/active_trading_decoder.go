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
		if a.u.MevSmartContractTxMapUniversalRouterOld.Txs == nil {
			a.u.MevSmartContractTxMapUniversalRouterOld.Txs = []web3_client.MevTx{}
		}
		methodName, args, err := web3_client.DecodeTxArgData(ctx, tx, a.u.MevSmartContractTxMapUniversalRouterOld)
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
		tmp := a.u.MevSmartContractTxMapUniversalRouterOld.Txs
		mevTxs = append(mevTxs, singleTx)
		a.u.MevSmartContractTxMapUniversalRouterOld.Txs = tmp
		a.m.TxFetcherMetrics.TransactionGroup(web3_client.UniswapUniversalRouterAddressOld, methodName)
	case web3_client.UniswapUniversalRouterAddressNew:
		if a.u.MevSmartContractTxMapUniversalRouterNew.Txs == nil {
			a.u.MevSmartContractTxMapUniversalRouterNew.Txs = []web3_client.MevTx{}
		}
		methodName, args, err := web3_client.DecodeTxArgData(ctx, tx, a.u.MevSmartContractTxMapUniversalRouterNew)
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
		tmp := a.u.MevSmartContractTxMapUniversalRouterNew.Txs
		mevTxs = append(mevTxs, singleTx)
		a.u.MevSmartContractTxMapUniversalRouterNew.Txs = tmp
		a.m.TxFetcherMetrics.TransactionGroup(web3_client.UniswapUniversalRouterAddressNew, methodName)
	case web3_client.UniswapV2Router02Address:
		if a.u.MevSmartContractTxMapV2Router02.Txs == nil {
			a.u.MevSmartContractTxMapV2Router02.Txs = []web3_client.MevTx{}
		}
		methodName, args, err := web3_client.DecodeTxArgData(ctx, tx, a.u.MevSmartContractTxMapV2Router02)
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
		tmp := a.u.MevSmartContractTxMapV2Router02.Txs
		mevTxs = append(mevTxs, singleTx)
		a.u.MevSmartContractTxMapV2Router02.Txs = tmp
		a.m.TxFetcherMetrics.TransactionGroup(web3_client.UniswapV2Router02Address, methodName)
	case web3_client.UniswapV2Router01Address:
		if a.u.MevSmartContractTxMapV2Router01.Txs == nil {
			a.u.MevSmartContractTxMapV2Router01.Txs = []web3_client.MevTx{}
		}
		methodName, args, err := web3_client.DecodeTxArgData(ctx, tx, a.u.MevSmartContractTxMapV2Router01)
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
		tmp := a.u.MevSmartContractTxMapV2Router01.Txs
		mevTxs = append(mevTxs, singleTx)
		a.u.MevSmartContractTxMapV2Router01.Txs = tmp
		a.m.TxFetcherMetrics.TransactionGroup(web3_client.UniswapV2Router01Address, methodName)
	case web3_client.UniswapV3Router01Address:
		if a.u.MevSmartContractTxMapV3SwapRouterV1.Txs == nil {
			a.u.MevSmartContractTxMapV3SwapRouterV1.Txs = []web3_client.MevTx{}
		}
		methodName, args, err := web3_client.DecodeTxArgData(ctx, tx, a.u.MevSmartContractTxMapV3SwapRouterV1)
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
		tmp := a.u.MevSmartContractTxMapV3SwapRouterV1.Txs
		mevTxs = append(mevTxs, singleTx)
		a.u.MevSmartContractTxMapV3SwapRouterV1.Txs = tmp
		a.m.TxFetcherMetrics.TransactionGroup(web3_client.UniswapV3Router01Address, methodName)
	case web3_client.UniswapV3Router02Address:
		if a.u.MevSmartContractTxMapV3SwapRouterV2.Txs == nil {
			a.u.MevSmartContractTxMapV3SwapRouterV2.Txs = []web3_client.MevTx{}
		}
		methodName, args, err := web3_client.DecodeTxArgData(ctx, tx, a.u.MevSmartContractTxMapV3SwapRouterV2)
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
		tmp := a.u.MevSmartContractTxMapV3SwapRouterV2.Txs
		mevTxs = append(mevTxs, singleTx)
		a.u.MevSmartContractTxMapV3SwapRouterV2.Txs = tmp
		a.m.TxFetcherMetrics.TransactionGroup(web3_client.UniswapV3Router02Address, methodName)
	}
	return mevTxs, nil
}
