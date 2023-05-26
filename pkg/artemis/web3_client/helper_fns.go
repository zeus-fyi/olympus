package web3_client

import (
	"context"
	"math/big"

	"github.com/gochain/gochain/v4/common"
)

func StringsToAddresses(addressOne, addressTwo string) (common.Address, common.Address) {
	addrOne := common.HexToAddress(addressOne)
	addrTwo := common.HexToAddress(addressTwo)
	return addrOne, addrTwo
}

type TxLifecycleStats struct {
	TxHash     common.Hash
	GasUsed    uint64
	TxBlockNum *big.Int
	RxBlockNum uint64
}

func (w *Web3Client) GetTxLifecycleStats(ctx context.Context, txHash common.Hash) (TxLifecycleStats, error) {
	tx, err := w.GetTransactionByHash(ctx, txHash)
	if err != nil {
		return TxLifecycleStats{}, err
	}
	rx, err := w.GetTransactionReceipt(ctx, txHash)
	if err != nil {
		return TxLifecycleStats{}, err
	}
	return TxLifecycleStats{
		TxHash:     txHash,
		GasUsed:    rx.GasUsed * tx.GasPrice.Uint64(),
		TxBlockNum: tx.BlockNumber,
		RxBlockNum: rx.BlockNumber,
	}, err
}
