package web3_client

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

func (w *Web3Client) DeployERC20Token(ctx context.Context, bc string, scParams web3_actions.SendContractTxPayload) (*types.Transaction, error) {
	w.Dial()
	defer w.Close()
	tx, err := w.DeployContract(ctx, bc, scParams)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (w *Web3Client) DeployRawdawgContract(ctx context.Context, bc string, scParams web3_actions.SendContractTxPayload) (*types.Transaction, error) {
	w.Dial()
	defer w.Close()
	tx, err := w.DeployContract(ctx, bc, scParams)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
