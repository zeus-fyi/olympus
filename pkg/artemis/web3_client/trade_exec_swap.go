package web3_client

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
	signing_automation_ethereum "github.com/zeus-fyi/zeus/pkg/artemis/signing_automation/ethereum"
)

func MustLoadSwapAbi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(artemis_oly_contract_abis.SwapABI))
	if err != nil {
		panic(err)
	}
	return readAbi
}

func LoadSwapAbiPayload(pairContractAddr string) (web3_actions.SendContractTxPayload, string, error) {
	abf := MustLoadSwapAbi()
	params := web3_actions.SendContractTxPayload{
		SmartContractAddr: pairContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractFile:      "",
		ContractABI:       abf,
		MethodName:        swap,
		Params:            []interface{}{},
	}
	return params, "", nil
}

func MustLoadTradingContractAbi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(artemis_oly_contract_abis.TradingAbi))
	if err != nil {
		panic(err)
	}
	return readAbi
}

func GetTradingSwapAbiPayload(tradingSwapContractAddr, pairContractAddr string, to *TradeOutcome, isToken0 bool) web3_actions.SendContractTxPayload {
	params := web3_actions.SendContractTxPayload{
		SmartContractAddr: tradingSwapContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       MustLoadTradingContractAbi(),
		MethodName:        execSmartContractTradingSwap,
		Params:            []interface{}{pairContractAddr, to.AmountInAddr.String(), to.AmountOutAddr.String(), to.AmountOut.String(), isToken0},
	}
	return params
}

func (u *UniswapV2Client) ExecSmartContractTradingSwap(pair UniswapV2Pair, to *TradeOutcome) (*web3_actions.SendContractTxPayload, error) {
	tokenNum := pair.GetTokenNumber(to.AmountInAddr)
	scInfo := GetTradingSwapAbiPayload("", pair.PairContractAddr, to, tokenNum == 0)

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
