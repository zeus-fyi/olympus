package web3_client

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

func (u *UniswapV2Client) RouterApproveAndSend(ctx context.Context, to *TradeOutcome, pairContractAddr string) error {
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
