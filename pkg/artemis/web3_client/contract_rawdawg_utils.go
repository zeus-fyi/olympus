package web3_client

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zeus-fyi/gochain/web3/accounts"
)

func (u *UniswapV2Client) ExecSmartContractTradingSwap(pair UniswapV2Pair, to *TradeOutcome) (*types.Transaction, error) {
	tokenNum := pair.GetTokenNumber(to.AmountInAddr)
	scInfo := GetRawdawgSwapAbiPayload(RawDawgAddr, pair.PairContractAddr, to, tokenNum == 0)

	// TODO implement better gas estimation
	scInfo.GasLimit = 3000000
	signedTx, err := u.Web3Client.GetSignedTxToCallFunctionWithArgs(ctx, &scInfo)
	if err != nil {
		return nil, err
	}
	to.AddTxHash(accounts.Hash(signedTx.Hash()))
	err = u.Web3Client.SendSignedTransaction(ctx, signedTx)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}
