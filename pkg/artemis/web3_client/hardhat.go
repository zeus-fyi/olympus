package web3_client

import (
	"context"
	"math/big"

	"github.com/zeus-fyi/gochain/v4/common/hexutil"
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

func (w *Web3Client) ResetNetwork(ctx context.Context, nodeURL string, blockNumber int) error {
	w.Dial()
	defer w.Close()
	err := w.Client.ResetNetwork(ctx, nodeURL, blockNumber)
	if err != nil {
		return err
	}
	return err
}

func (w *Web3Client) ImpersonateAccount(ctx context.Context, userAddr string) error {
	w.Dial()
	defer w.Close()
	err := w.Client.ImpersonateAccount(ctx, userAddr)
	if err != nil {
		return err
	}
	return err
}

func (w *Web3Client) SetStorageAt(ctx context.Context, addr, slot, value string) error {
	w.Dial()
	defer w.Close()
	err := w.Client.SetStorageAt(ctx, addr, slot, value)
	if err != nil {
		return err
	}
	return err
}

func (w *Web3Client) GetStorageAt(ctx context.Context, addr, slot string) (hexutil.Bytes, error) {
	w.Dial()
	defer w.Close()
	result, err := w.Client.GetStorageAt(ctx, addr, slot)
	if err != nil {
		return result, err
	}
	return result, err
}

func (w *Web3Client) GetEvmSnapshot(ctx context.Context) (*big.Int, error) {
	w.Dial()
	defer w.Close()
	ss, err := w.Client.GetEVMSnapshot(ctx)
	if err != nil {
		return ss, err
	}
	return ss, err
}
