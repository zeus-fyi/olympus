package artemis_rawdawg_contract

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	artemis_utils "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/utils"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	web3_actions "github.com/zeus-fyi/zeus/pkg/artemis/web3/client"
)

type RawDawgV2SimSwapParams struct {
	Pair      common.Address `json:"_pair"`
	TokenIn   common.Address `json:"_token_in"`
	TokenOut  common.Address `json:"_token_out"`
	AmountIn  *big.Int       `json:"_amountIn"`
	AmountOut *big.Int       `json:"_amountOut"`
	IsToken0  bool           `json:"_isToken0"`
}

func GetRawdawgV2SimSwapAbiPayload(tradingSwapContractAddr string, to *artemis_trading_types.TradeOutcome) *web3_actions.SendContractTxPayload {
	isToken0 := false
	pairContractAddr, tkn0, _ := artemis_utils.CreateV2TradingPair(to.AmountInAddr, to.AmountOutAddr)
	if tkn0.String() == to.AmountInAddr.String() {
		isToken0 = true
	}
	params := &web3_actions.SendContractTxPayload{
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
