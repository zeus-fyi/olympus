package web3_client

import (
	"context"

	"github.com/gochain/gochain/v4/common/hexutil"
)

func (w *Web3Client) SetCodeOverride(ctx context.Context, addr, byteCode string) error {
	w.Dial()
	defer w.Close()
	err := w.SetCode(ctx, addr, hexutil.Bytes(byteCode))
	if err != nil {
		return err
	}
	return err
}
