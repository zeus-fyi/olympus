package web3_client

import (
	"context"
	"math/big"

	web3_types "github.com/zeus-fyi/gochain/web3/types"
)

func (w *Web3Client) GetCurrentBalance(ctx context.Context) (*big.Int, error) {
	w.Dial()
	defer w.Close()
	b, err := w.GetBalance(ctx, w.PublicKey(), nil)
	if err != nil {
		return nil, err
	}
	return b, err
}

func (w *Web3Client) GetUserCurrentBalance(ctx context.Context, userAddr string) (*big.Int, error) {
	w.Dial()
	defer w.Close()
	b, err := w.GetBalance(ctx, userAddr, nil)
	if err != nil {
		return nil, err
	}
	return b, err
}

func (w *Web3Client) GetCurrentBalanceGwei(ctx context.Context) (string, error) {
	w.Dial()
	defer w.Close()
	b, err := w.GetBalance(ctx, w.PublicKey(), nil)
	if err != nil {
		return "0", err
	}
	return web3_types.WeiAsGwei(b), err
}
