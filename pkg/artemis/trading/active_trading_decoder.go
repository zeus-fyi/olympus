package artemis_realtime_trading

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
)

func (a *ActiveTrading) DecodeTx(ctx context.Context, tx *types.Transaction) {

	// todo
}

/*
func (u *UniswapClient) ProcessMempoolTxs(ctx context.Context, mempool map[string]map[string]*types.Transaction) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	for userAddr, txPoolQueue := range mempool {
		for order, tx := range txPoolQueue {
			if tx.To() == nil {
				continue
			}
			switch tx.To().String() {
			case UniswapUniversalRouterAddressOld:
				if u.MevSmartContractTxMapUniversalRouterOld.Txs == nil {
					u.MevSmartContractTxMapUniversalRouterOld.Txs = []MevTx{}
				}
				methodName, args, err := DecodeTxArgData(ctx, tx, u.MevSmartContractTxMapUniversalRouterOld)
				if err != nil {
					continue
				}
				singleTx := MevTx{
					MethodName:  methodName,
					UserAddr:    userAddr,
					Args:        args,
					Order:       order,
					TxPoolQueue: txPoolQueue,
					Tx:          tx,
				}

				tmp := u.MevSmartContractTxMapUniversalRouterOld.Txs
				tmp = append(tmp, singleTx)
				u.MevSmartContractTxMapUniversalRouterOld.Txs = tmp
			case UniswapUniversalRouterAddressNew:
				if u.MevSmartContractTxMapUniversalRouterNew.Txs == nil {
					u.MevSmartContractTxMapUniversalRouterNew.Txs = []MevTx{}
				}
				methodName, args, err := DecodeTxArgData(ctx, tx, u.MevSmartContractTxMapUniversalRouterNew)
				if err != nil {
					continue
				}
				singleTx := MevTx{
					MethodName:  methodName,
					UserAddr:    userAddr,
					Args:        args,
					Order:       order,
					TxPoolQueue: txPoolQueue,
					Tx:          tx,
				}

				tmp := u.MevSmartContractTxMapUniversalRouterNew.Txs
				tmp = append(tmp, singleTx)
				u.MevSmartContractTxMapUniversalRouterNew.Txs = tmp
			case UniswapV2Router02Address:
				if u.MevSmartContractTxMapV2Router02.Txs == nil {
					u.MevSmartContractTxMapV2Router02.Txs = []MevTx{}
				}
				methodName, args, err := DecodeTxArgData(ctx, tx, u.MevSmartContractTxMapV2Router02)
				if err != nil {
					continue
				}
				singleTx := MevTx{
					MethodName:  methodName,
					UserAddr:    userAddr,
					Args:        args,
					Order:       order,
					TxPoolQueue: txPoolQueue,
					Tx:          tx,
				}
				tmp := u.MevSmartContractTxMapV2Router02.Txs
				tmp = append(tmp, singleTx)
				u.MevSmartContractTxMapV2Router02.Txs = tmp
			case UniswapV2Router01Address:
				if u.MevSmartContractTxMapV2Router01.Txs == nil {
					u.MevSmartContractTxMapV2Router01.Txs = []MevTx{}
				}
				methodName, args, err := DecodeTxArgData(ctx, tx, u.MevSmartContractTxMapV2Router01)
				if err != nil {
					continue
				}
				singleTx := MevTx{
					MethodName:  methodName,
					UserAddr:    userAddr,
					Args:        args,
					Order:       order,
					TxPoolQueue: txPoolQueue,
					Tx:          tx,
				}
				tmp := u.MevSmartContractTxMapV2Router01.Txs
				tmp = append(tmp, singleTx)
				u.MevSmartContractTxMapV2Router01.Txs = tmp
			case UniswapV3Router01Address:
				if u.MevSmartContractTxMapV3SwapRouterV1.Txs == nil {
					u.MevSmartContractTxMapV3SwapRouterV1.Txs = []MevTx{}
				}
				methodName, args, err := DecodeTxArgData(ctx, tx, u.MevSmartContractTxMapV3SwapRouterV1)
				if err != nil {
					continue
				}
				singleTx := MevTx{
					MethodName:  methodName,
					UserAddr:    userAddr,
					Args:        args,
					Order:       order,
					TxPoolQueue: txPoolQueue,
					Tx:          tx,
				}
				tmp := u.MevSmartContractTxMapV3SwapRouterV1.Txs
				tmp = append(tmp, singleTx)
				u.MevSmartContractTxMapV3SwapRouterV1.Txs = tmp
			case UniswapV3Router02Address:
				if u.MevSmartContractTxMapV3SwapRouterV2.Txs == nil {
					u.MevSmartContractTxMapV3SwapRouterV2.Txs = []MevTx{}
				}
				methodName, args, err := DecodeTxArgData(ctx, tx, u.MevSmartContractTxMapV3SwapRouterV2)
				if err != nil {
					continue
				}
				singleTx := MevTx{
					MethodName:  methodName,
					UserAddr:    userAddr,
					Args:        args,
					Order:       order,
					TxPoolQueue: txPoolQueue,
					Tx:          tx,
				}
				tmp := u.MevSmartContractTxMapV3SwapRouterV2.Txs
				tmp = append(tmp, singleTx)
				u.MevSmartContractTxMapV3SwapRouterV2.Txs = tmp
			}
		}
	}
	return nil
}

*/
