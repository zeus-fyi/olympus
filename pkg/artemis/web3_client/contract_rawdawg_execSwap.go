package web3_client

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

type RawDawgSwapParams struct {
	Pair      common.Address `json:"_pair"`
	TokenIn   common.Address `json:"_token_in"`
	AmountIn  *big.Int       `json:"_amountIn"`
	AmountOut *big.Int       `json:"_amountOut"`
	IsToken0  bool           `json:"_isToken0"`
}

func GetRawdawgSwapAbiPayload(tradingSwapContractAddr, pairContractAddr string, to *TradeOutcome, isToken0 bool) web3_actions.SendContractTxPayload {
	params := web3_actions.SendContractTxPayload{
		SmartContractAddr: tradingSwapContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       RawdawgAbi,
		MethodName:        execSmartContractTradingSwap,
		Params:            []interface{}{pairContractAddr, to.AmountInAddr.String(), isToken0, to.AmountIn.String(), to.AmountOut.String()},
	}
	return params
}

func (u *UniswapClient) ExecSmartContractTradingSwap(tradingContractAddr string, pair UniswapV2Pair, to *TradeOutcome) (*types.Transaction, error) {
	tokenNum := pair.GetTokenNumber(to.AmountInAddr)
	scInfo := GetRawdawgSwapAbiPayload(tradingContractAddr, pair.PairContractAddr, to, tokenNum == 0)
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
