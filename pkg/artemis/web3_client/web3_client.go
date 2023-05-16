package web3_client

import (
	"context"
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/gochain/web3/web3_actions"
)

type Web3Client struct {
	web3_actions.Web3Actions
}

func NewWeb3Client(nodeUrl string, acc *accounts.Account) Web3Client {
	w := web3_actions.NewWeb3ActionsClientWithAccount(nodeUrl, acc)
	return Web3Client{w}
}

func (w *Web3Client) GetBlockHeight(ctx context.Context) (*big.Int, error) {
	w.Dial()
	defer w.Close()
	blockNumber, err := w.GetBlockNumber(ctx)
	if err != nil {
		return nil, err
	}
	return blockNumber, nil
}
