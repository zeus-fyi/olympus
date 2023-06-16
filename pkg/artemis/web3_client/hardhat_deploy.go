package web3_client

import (
	"context"
)

func (w *Web3Client) SetCodeOverride(ctx context.Context, addr, byteCode string) error {
	w.Dial()
	defer w.Close()
	err := w.SetCode(ctx, addr, byteCode)
	if err != nil {
		return err
	}
	return err
}
