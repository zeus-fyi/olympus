package web3_client

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pingcap/errors"
)

func (w *Web3Client) HardHatSetBalance(ctx context.Context, addr string, balance hexutil.Big) error {
	w.Dial()
	defer w.Close()
	err := w.SetBalance(ctx, addr, balance)
	if err != nil {
		return err
	}
	return err
}

const irisBetaSvc = "https://iris.zeus.fyi/v1beta/internal/"

func (w *Web3Client) HardHatResetNetwork(ctx context.Context, blockNumber int) error {
	w.Dial()
	defer w.Close()

	nodeURL := ""
	nodeInfo, err := w.GetNodeInfo(ctx)
	if err != nil {
		return err
	}
	if nodeInfo.ForkConfig.ForkUrl != "" {
		nodeURL = nodeInfo.ForkConfig.ForkUrl
	}
	if nodeURL == irisBetaSvc {
		return errors.New("iris beta cannot proxy itself recursively")
	}
	if nodeURL == "" {
		return errors.New("node url is empty")
	}

	err = w.ResetNetwork(ctx, nodeURL, blockNumber)
	if err != nil {
		return err
	}
	return err
}

func (w *Web3Client) HardhatImpersonateAccount(ctx context.Context, userAddr string) error {
	w.Dial()
	defer w.Close()
	err := w.ImpersonateAccount(ctx, userAddr)
	if err != nil {
		return err
	}
	return err
}

func (w *Web3Client) HardhatSetStorageAt(ctx context.Context, addr, slot, value string) error {
	w.Dial()
	defer w.Close()
	err := w.SetStorageAt(ctx, addr, slot, value)
	if err != nil {
		return err
	}
	return err
}

func (w *Web3Client) HardHatGetStorageAt(ctx context.Context, addr, slot string) (hexutil.Bytes, error) {
	w.Dial()
	defer w.Close()
	result, err := w.GetStorageAt(ctx, addr, slot)
	if err != nil {
		return result, err
	}
	return result, err
}

func (w *Web3Client) HardHatGetEvmSnapshot(ctx context.Context) (*big.Int, error) {
	w.Dial()
	defer w.Close()
	ss, err := w.GetEVMSnapshot(ctx)
	if err != nil {
		return ss, err
	}
	return ss, err
}

func (w *Web3Client) SetCodeOverride(ctx context.Context, addr, byteCode string) error {
	w.Dial()
	defer w.Close()
	err := w.SetCode(ctx, addr, byteCode)
	if err != nil {
		return err
	}
	return err
}
