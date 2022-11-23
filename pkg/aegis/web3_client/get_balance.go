package web3_client

import (
	"context"
	"math/big"

	"github.com/gochain/web3"
)

func (w *Web3Client) GetCurrentBalance(ctx context.Context) (*big.Int, error) {
	w.Dial()
	defer w.Close()
	b, err := w.Client.GetBalance(ctx, w.Address().Hex(), nil)
	return b, err
}

func (w *Web3Client) GetCurrentBalanceGwei(ctx context.Context) (string, error) {
	w.Dial()
	defer w.Close()
	b, err := w.Client.GetBalance(ctx, w.Address().Hex(), nil)
	if err != nil {
		return "0", err
	}
	return web3.WeiAsGwei(b), err
}
