package web3_client

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

func (w *Web3Client) ValidateTxIsPending(ctx context.Context, txHashStr string) (bool, error) {
	w.Dial()
	defer w.Close()
	txHash := common.HexToHash(txHashStr)
	_, isPending, err := w.Web3Actions.C.TransactionByHash(ctx, txHash)
	if err != nil {
		return false, err
	}
	if isPending {
		return true, nil
	}
	return false, nil
}
