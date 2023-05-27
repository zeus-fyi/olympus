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

func (w *Web3Client) GetEthBalance(ctx context.Context, addr string, blockNum *big.Int) (*big.Int, error) {
	w.Dial()
	defer w.Close()
	balance, err := w.GetBalance(ctx, addr, blockNum)
	if err != nil {
		return balance, err
	}
	return balance, err
}

func ConvertAmountsToBigIntSlice(amounts []interface{}) []*big.Int {
	var amountsBigInt []*big.Int
	for _, amount := range amounts {
		pair := amount.([]*big.Int)
		for _, p := range pair {
			amountsBigInt = append(amountsBigInt, p)
		}
	}
	return amountsBigInt
}
