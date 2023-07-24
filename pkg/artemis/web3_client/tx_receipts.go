package web3_client

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (w *Web3Client) GetTxReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, bool, error) {
	w.Dial()
	defer w.Close()
	rx, err := w.C.TransactionReceipt(ctx, txHash)
	if err != nil {
		if err.Error() == "not found" {
			return nil, false, nil
		}
		return nil, false, err
	}
	return rx, true, nil
}
