package web3_client

import (
	"context"
	"fmt"
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

type Web3Client struct {
	web3_actions.Web3Actions
}

func NewWeb3ClientFakeSigner(nodeUrl string) Web3Client {
	acc, _ := accounts.CreateAccount()
	w := web3_actions.NewWeb3ActionsClientWithAccount(nodeUrl, acc)
	return Web3Client{w}
}

func NewWeb3ClientWithRelay(nodeUrl, relayUrl string, acc *accounts.Account) Web3Client {
	w := web3_actions.NewWeb3ActionsClientWithRelay(nodeUrl, relayUrl, acc)
	return Web3Client{w}
}

func NewWeb3Client(nodeUrl string, acc *accounts.Account) Web3Client {
	w := web3_actions.NewWeb3ActionsClientWithAccount(nodeUrl, acc)
	return Web3Client{w}
}

func (w *Web3Client) GetBlockHeight(ctx context.Context) (*big.Int, error) {
	w.Dial()
	defer w.Close()
	blockNumber, err := w.C.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetUint64(blockNumber), nil
}

func (w *Web3Client) GetHeadBlockHeight(ctx context.Context) (*big.Int, error) {
	w.Dial()
	defer w.C.Close()
	isSyncing, err := w.C.SyncProgress(ctx)
	if err != nil {
		fmt.Println("error getting sync status", isSyncing)
		return nil, err
	}
	if isSyncing != nil && isSyncing.CurrentBlock != isSyncing.HighestBlock {
		return nil, fmt.Errorf("node is not synced")
	}
	return w.GetBlockHeight(ctx)
}
