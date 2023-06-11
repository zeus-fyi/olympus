package web3_client

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

const (
	batchExecuteSwap = "batchExecuteSwap"
)

type BatchRawDawgParams struct {
	Swap []RawDawgSwapParams `json:"_swap"`
}

func (r *BatchRawDawgParams) AddRawdawgSwapParams(pair UniswapV2Pair, to *TradeOutcome) {
	if r.Swap == nil {
		r.Swap = []RawDawgSwapParams{}
	}
	isToken0 := false
	tokenNum := pair.GetTokenNumber(to.AmountInAddr)
	if tokenNum == 0 {
		isToken0 = true
	}
	params := RawDawgSwapParams{
		Pair:      common.HexToAddress(pair.PairContractAddr),
		TokenIn:   common.HexToAddress(to.AmountInAddr.Hex()),
		AmountIn:  to.AmountIn,
		AmountOut: to.AmountOut,
		IsToken0:  isToken0,
	}
	r.Swap = append(r.Swap, params)
}

func GetRawdawgSwapAbiBatchPayload(tradingSwapContractAddr string, rawDawgBatch BatchRawDawgParams) web3_actions.SendContractTxPayload {
	params := web3_actions.SendContractTxPayload{
		SmartContractAddr: tradingSwapContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       RawdawgAbi,
		MethodName:        batchExecuteSwap,
		Params:            []interface{}{rawDawgBatch.Swap},
	}
	return params
}

func (u *UniswapClient) ExecSmartContractTradingBatchSwap(tradingContractAddr string, params BatchRawDawgParams) (*types.Transaction, error) {
	scInfo := GetRawdawgSwapAbiBatchPayload(tradingContractAddr, params)
	scInfo.MethodName = batchExecuteSwap
	// TODO implement better gas estimation
	scInfo.GasLimit = 3000000
	signedTx, err := u.Web3Client.GetSignedTxToCallFunctionWithArgs(ctx, &scInfo)
	if err != nil {
		return nil, err
	}
	err = u.Web3Client.SendSignedTransaction(ctx, signedTx)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}
