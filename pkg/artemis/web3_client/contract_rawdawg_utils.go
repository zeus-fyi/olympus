package web3_client

import (
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

func (u *UniswapV2Client) ExecSmartContractTradingSwap(pair UniswapV2Pair, to *TradeOutcome) (*web3_actions.SendContractTxPayload, error) {
	tokenNum := pair.GetTokenNumber(to.AmountInAddr)
	scInfo := GetRawdawgSwapAbiPayload(RawDawgAddr, pair.PairContractAddr, to, tokenNum == 0)

	// TODO implement better gas estimation
	scInfo.GasLimit = 3000000
	signedTx, err := u.Web3Client.GetSignedTxToCallFunctionWithArgs(ctx, &scInfo)
	if err != nil {
		return &web3_actions.SendContractTxPayload{}, err
	}
	to.AddTxHash(accounts.Hash(signedTx.Hash()))
	err = u.Web3Client.SendSignedTransaction(ctx, signedTx)
	if err != nil {
		return &web3_actions.SendContractTxPayload{}, err
	}
	return &scInfo, nil
}
