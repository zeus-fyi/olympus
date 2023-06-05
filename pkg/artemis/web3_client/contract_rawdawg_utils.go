package web3_client

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

const (
	batchExecuteSwap = "batchExecuteSwap"
)

type RawDawgParams struct {
	Pair      common.Address `json:"_pair"`
	TokenIn   common.Address `json:"_token_in"`
	AmountIn  *big.Int       `json:"_amountIn"`
	AmountOut *big.Int       `json:"_amountOut"`
	IsToken0  bool           `json:"_isToken0"`
}

type BatchRawDawgParams struct {
	Swap []RawDawgParams `json:"_swap"`
}

func (r *BatchRawDawgParams) AddRawdawgParams(pair UniswapV2Pair, to *TradeOutcome) {
	if r.Swap == nil {
		r.Swap = []RawDawgParams{}
	}
	isToken0 := false
	tokenNum := pair.GetTokenNumber(to.AmountInAddr)
	if tokenNum == 0 {
		isToken0 = true
	}
	params := RawDawgParams{
		Pair:      common.HexToAddress(pair.PairContractAddr),
		TokenIn:   common.HexToAddress(to.AmountInAddr.Hex()),
		AmountIn:  to.AmountIn,
		AmountOut: to.AmountOut,
		IsToken0:  isToken0,
	}
	r.Swap = append(r.Swap, params)
}

func GetRawdawgSwapAbiPayload(tradingSwapContractAddr, pairContractAddr string, to *TradeOutcome, isToken0 bool) web3_actions.SendContractTxPayload {
	params := web3_actions.SendContractTxPayload{
		SmartContractAddr: tradingSwapContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       MustLoadRawdawgAbi(),
		MethodName:        execSmartContractTradingSwap,
		Params:            []interface{}{pairContractAddr, to.AmountInAddr.String(), to.AmountIn.String(), to.AmountOut.String(), isToken0},
	}
	return params
}

/*
	out, err := swapAbi.Methods[batchExecuteSwap].Inputs.Pack(rawDawgBatch.Swap)
	if err != nil {
		panic(err)
	}
	output, err := swapAbi.Methods[batchExecuteSwap].Inputs.Unpack(out)
	if err != nil {
		panic(err)
	}
	fmt.Println("paramsInput", output)
*/

func GetRawdawgSwapAbiBatchPayload(tradingSwapContractAddr string, rawDawgBatch BatchRawDawgParams) web3_actions.SendContractTxPayload {
	swapAbi := MustLoadRawdawgAbi()
	params := web3_actions.SendContractTxPayload{
		SmartContractAddr: tradingSwapContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       swapAbi,
		MethodName:        batchExecuteSwap,
		Params:            []interface{}{rawDawgBatch.Swap},
	}
	return params
}

func (u *UniswapV2Client) ExecSmartContractTradingBatchSwap(tradingContractAddr string, params BatchRawDawgParams) (*types.Transaction, error) {
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

func (u *UniswapV2Client) ExecSmartContractTradingSwap(tradingContractAddr string, pair UniswapV2Pair, to *TradeOutcome) (*types.Transaction, error) {
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
