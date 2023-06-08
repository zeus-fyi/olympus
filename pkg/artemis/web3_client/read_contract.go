package web3_client

import (
	"context"

	web3_types "github.com/zeus-fyi/gochain/web3/types"
)

func (w *Web3Client) ReadContract(ctx context.Context, abiFile, address string) (string, error) {
	w.Dial()
	defer w.Close()
	b, err := w.GetBalance(ctx, w.PublicKey(), nil)
	if err != nil {
		return "0", err
	}
	return web3_types.WeiAsGwei(b), err
}
