package web3_client

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/v4/common"
	"github.com/zeus-fyi/gochain/web3/web3_actions"
)

func (u *UniswapV2Client) RouterApproveAndSend(ctx context.Context, to *TradeOutcome, pairContractAddr string) error {
	approveTx, err := u.Web3Client.ERC20ApproveSpender(ctx, to.AmountInAddr.String(), u.RouterSmartContractAddr, to.AmountIn)
	if err != nil {
		log.Warn().Interface("approveTx", approveTx).Err(err).Msg("error approving router")
		return err
	}
	to.AddTxHash(approveTx.Hash)
	transferTxParams := web3_actions.SendContractTxPayload{
		SmartContractAddr: to.AmountInAddr.String(),
		SendEtherPayload: web3_actions.SendEtherPayload{
			TransferArgs: web3_actions.TransferArgs{
				ToAddress: common.HexToAddress(pairContractAddr),
			},
		},
		ContractABI: LoadERC20Abi(),
		Params:      []interface{}{common.HexToAddress(pairContractAddr), to.AmountIn},
	}
	transferTx, err := u.Web3Client.TransferERC20Token(ctx, transferTxParams)
	if err != nil {
		log.Warn().Interface("transferTx", transferTx).Err(err).Msg("error approving router")
		return err
	}
	to.AddTxHash(transferTx.Hash)
	return err
}
