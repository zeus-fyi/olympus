package web3_client

import (
	"errors"
	"math/big"
	"strings"

	"github.com/zeus-fyi/gochain/web3/web3_actions"
)

func (u *UniswapV2Client) GetAmounts(to TradeOutcome, method string) ([]interface{}, error) {
	switch method {
	case getAmountsOut:
		pathSlice := []string{to.AmountInAddr.String(), to.AmountOutAddr.String()}
		return u.GetAmountsOut(to.AmountIn, pathSlice)
	case getAmountsIn:
		pathSlice := []string{to.AmountOutAddr.String(), to.AmountInAddr.String()}
		return u.GetAmountsIn(to.AmountOut, pathSlice)
	}
	return nil, errors.New("invalid method")
}

/*
Given an output asset amount and an array of token addresses, calculates all preceding minimum
input token amounts by calling getReserves for each pair of token addresses in the path in turn,
and using these to call getAmountIn.
*/
func (u *UniswapV2Client) GetAmountsIn(amountOut *big.Int, pathSlice []string) ([]interface{}, error) {
	pathString := "[" + strings.Join(pathSlice, ",") + "]"
	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: u.MevSmartContractTxMap.SmartContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       u.MevSmartContractTxMap.Abi,
		MethodName:        getAmountsIn,
		Params:            []interface{}{amountOut, pathString},
	}
	amountsIn, err := u.Web3Client.GetContractConst(ctx, scInfo)
	if err != nil {
		return nil, err
	}
	return amountsIn, err
}

/*
Given an output asset amount and an array of token addresses, calculates all preceding minimum
input token amounts by calling getReserves for each pair of token addresses in the path in turn,
and using these to call getAmountIn.
*/
func (u *UniswapV2Client) GetAmountsOut(amountIn *big.Int, pathSlice []string) ([]interface{}, error) {
	pathString := "[" + strings.Join(pathSlice, ",") + "]"
	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: u.MevSmartContractTxMap.SmartContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       u.MevSmartContractTxMap.Abi,
		MethodName:        getAmountsOut,
		Params:            []interface{}{amountIn, pathString},
	}
	amountsOut, err := u.Web3Client.GetContractConst(ctx, scInfo)
	if err != nil {
		return nil, err
	}
	return amountsOut, err
}
