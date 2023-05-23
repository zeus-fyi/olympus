package web3_client

import (
	"context"

	"github.com/gochain/gochain/v4/common/hexutil"
)

func (w *Web3Client) SetBalance(ctx context.Context, addr string, balance hexutil.Big) error {
	w.Dial()
	defer w.Close()
	err := w.Client.SetBalance(ctx, addr, balance)
	if err != nil {
		return err
	}
	return err
}
