package web3_client

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

func (w *Web3Client) GetBlockTxs(ctx context.Context) (types.Transactions, error) {
	w.Dial()
	defer w.Close()
	block, err := w.C.BlockByNumber(ctx, nil)
	if err != nil {
		log.Err(err).Msg("failed to get nonce")
		return nil, err
	}
	return block.Transactions(), nil
}

func (w *Web3Client) GetTxByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
	w.Dial()
	defer w.Close()
	tx, isPending, err := w.C.TransactionByHash(ctx, hash)
	if err != nil {
		log.Err(err).Msg("failed to get nonce")
		return nil, false, err
	}
	return tx, isPending, nil
}

func (u *UniswapClient) RouterApproveAndSend(ctx context.Context, to *TradeOutcome, pairContractAddr string) error {
	approveTx, err := u.Web3Client.ERC20ApproveSpender(ctx, to.AmountInAddr.String(), u.RouterSmartContractAddr, to.AmountIn)
	if err != nil {
		log.Warn().Interface("approveTx", approveTx).Err(err).Msg("error approving router")
		return err
	}
	to.AddTxHash(accounts.Hash(approveTx.Hash()))
	transferTxParams := web3_actions.SendContractTxPayload{
		SmartContractAddr: to.AmountInAddr.String(),
		SendEtherPayload: web3_actions.SendEtherPayload{
			TransferArgs: web3_actions.TransferArgs{
				ToAddress: accounts.HexToAddress(pairContractAddr),
			},
		},
		ContractABI: MustLoadERC20Abi(),
		Params:      []interface{}{accounts.HexToAddress(pairContractAddr), to.AmountIn},
	}
	transferTx, err := u.Web3Client.TransferERC20Token(ctx, transferTxParams)
	if err != nil {
		log.Warn().Interface("transferTx", transferTx).Err(err).Msg("error approving router")
		return err
	}
	to.AddTxHash(accounts.Hash(transferTx.Hash()))
	return err
}
