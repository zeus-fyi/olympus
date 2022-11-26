package web3_client

import (
	"context"

	"github.com/gochain/web3"
)

func (w *Web3Client) ReadContract(ctx context.Context, abiFile, address string) (string, error) {
	w.Dial()
	defer w.Close()
	b, err := w.Client.GetBalance(ctx, w.Address().Hex(), nil)
	if err != nil {
		return "0", err
	}
	return web3.WeiAsGwei(b), err
}
