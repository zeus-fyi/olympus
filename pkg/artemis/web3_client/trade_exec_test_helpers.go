package web3_client

import (
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
)

func (u *UniswapClient) ExecSwap(pair uniswap_pricing.UniswapV2Pair, to *artemis_trading_types.TradeOutcome) (*web3_actions.SendContractTxPayload, error) {
	scInfo, _, err := artemis_oly_contract_abis.LoadSwapAbiPayload(pair.PairContractAddr)
	if err != nil {
		return &web3_actions.SendContractTxPayload{}, err
	}
	tokenNum := pair.GetTokenNumber(to.AmountInAddr)
	if tokenNum == 0 {
		scInfo.Params = []interface{}{"0", to.AmountOut, u.Web3Client.Address(), []byte{}}
	} else {
		scInfo.Params = []interface{}{to.AmountOut, "0", u.Web3Client.Address(), []byte{}}
	}
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
