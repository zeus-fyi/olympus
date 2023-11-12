package artemis_rawdawg_contract

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	web3_actions "github.com/zeus-fyi/zeus/pkg/artemis/web3/client"
)

type RawDawgV2SimSwapParams struct {
	Pair      common.Address `json:"_pair"`
	TokenIn   common.Address `json:"_token_in"`
	AmountIn  *big.Int       `json:"_amountIn"`
	AmountOut *big.Int       `json:"_amountOut"`
	IsToken0  bool           `json:"_isToken0"`
}

func GetRawdawgV2SimSwapAbiPayload(tradingSwapContractAddr, pairContractAddr string, to *artemis_trading_types.TradeOutcome, isToken0 bool) web3_actions.SendContractTxPayload {
	params := web3_actions.SendContractTxPayload{
		SmartContractAddr: tradingSwapContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       RawdawgAbi,
		MethodName:        simulateV2AndRevertSwap,
		Params:            []interface{}{pairContractAddr, to.AmountInAddr.String(), to.AmountOutAddr.String(), isToken0, to.AmountIn.String(), to.AmountOut.String()},
	}
	return params
}

func ExecSmartContractTradingBatchSwap(w3a web3_actions.Web3Actions, tradingContractAddr string, params RawDawgV2SimSwapParams) (*types.Transaction, error) {
	return nil, nil
}

/*
ExecSmartContractTradingBatchSwap(tradingContractAddr string, params BatchRawDawgParams) (*types.Transaction, error) {
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
*/
