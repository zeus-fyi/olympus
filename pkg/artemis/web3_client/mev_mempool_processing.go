package web3_client

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

func (u *UniswapClient) ProcessMempoolTxs(ctx context.Context, mempool map[string]map[string]*types.Transaction) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	for userAddr, txPoolQueue := range mempool {
		for order, tx := range txPoolQueue {
			if tx.To() == nil {
				continue
			}
			switch tx.To().String() {
			case UniswapUniversalRouterAddressOld, UniswapUniversalRouterAddress:
				if u.MevSmartContractTxMapUniversalRouter.Txs == nil {
					u.MevSmartContractTxMapUniversalRouter.Txs = []MevTx{}
				}
				methodName, args, err := DecodeTxArgData(ctx, tx, u.MevSmartContractTxMapUniversalRouter)
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

				tmp := u.MevSmartContractTxMapUniversalRouter.Txs
				tmp = append(tmp, singleTx)
				u.MevSmartContractTxMapUniversalRouter.Txs = tmp
			case UniswapV2RouterAddress, UniswapV2RouterAddress2:
				if u.MevSmartContractTxMap.Txs == nil {
					u.MevSmartContractTxMap.Txs = []MevTx{}
				}
				methodName, args, err := DecodeTxArgData(ctx, tx, u.MevSmartContractTxMap)
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
				tmp := u.MevSmartContractTxMap.Txs
				tmp = append(tmp, singleTx)
				u.MevSmartContractTxMap.Txs = tmp
			}
		}
	}
	return nil
}
