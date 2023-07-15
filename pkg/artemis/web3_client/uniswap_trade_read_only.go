package web3_client

import (
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
)

func (u *UniswapClient) GetAmounts(address *common.Address, to artemis_trading_types.TradeOutcome, method string) ([]*big.Int, error) {
	switch method {
	case getAmountsOut:
		pathSlice := []string{to.AmountInAddr.String(), to.AmountOutAddr.String()}
		return u.GetAmountsOut(address, to.AmountIn, pathSlice)
	case getAmountsIn:
		pathSlice := []string{to.AmountOutAddr.String(), to.AmountInAddr.String()}
		return u.GetAmountsIn(address, to.AmountOut, pathSlice)
	}
	return nil, errors.New("invalid method")
}

/*
	Given an output asset amount and an array of token addresses, calculates all preceding minimum
	input token amounts by calling getReserves for each pair of token addresses in the path in turn,
	and using these to call getAmountIn.
*/

func (u *UniswapClient) GetAmountsIn(address *common.Address, amountOut *big.Int, pathSlice []string) ([]*big.Int, error) {
	mm := u.MevSmartContractTxMapV2Router02
	if address != nil {
		if address.String() == u.MevSmartContractTxMapV2Router01.SmartContractAddr {
			mm = u.MevSmartContractTxMapV2Router01
		}
	}
	pathString := "[" + strings.Join(pathSlice, ",") + "]"
	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: mm.SmartContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       mm.Abi,
		MethodName:        getAmountsIn,
		Params:            []interface{}{amountOut, pathString},
	}
	amountsIn, err := u.Web3Client.GetContractConst(ctx, scInfo)
	if err != nil {
		return nil, err
	}
	amountsInFirstPair := ConvertAmountsToBigIntSlice(amountsIn)
	return amountsInFirstPair, err
}

func (u *UniswapClient) GetAmountsOut(address *common.Address, amountIn *big.Int, pathSlice []string) ([]*big.Int, error) {
	mm := u.MevSmartContractTxMapV2Router02
	if address != nil {
		if address.String() == u.MevSmartContractTxMapV2Router01.SmartContractAddr {
			mm = u.MevSmartContractTxMapV2Router01
		}
	}
	pathString := "[" + strings.Join(pathSlice, ",") + "]"
	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: mm.SmartContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       mm.Abi,
		MethodName:        getAmountsOut,
		Params:            []interface{}{amountIn, pathString},
	}
	amountsOut, err := u.Web3Client.GetContractConst(ctx, scInfo)
	if err != nil {
		return nil, err
	}
	amountsOutFirstPair := ConvertAmountsToBigIntSlice(amountsOut)
	return amountsOutFirstPair, err
}
