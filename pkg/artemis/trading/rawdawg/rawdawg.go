package artemis_rawdawg_contract

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	artemis_utils "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/utils"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
	web3_actions "github.com/zeus-fyi/zeus/pkg/artemis/web3/client"
)

var (
	RawdawgAbi = artemis_oly_contract_abis.MustLoadRawdawgAbi()
)

type RawDawgSwapParams struct {
	Pair      common.Address `json:"_pair"`
	TokenIn   common.Address `json:"_token_in"`
	AmountIn  *big.Int       `json:"_amountIn"`
	AmountOut *big.Int       `json:"_amountOut"`
	IsToken0  bool           `json:"_isToken0"`
}

const (
	execSmartContractTradingSwap = "executeSwap"
)

func GetRawDawgSwapAbiPayload(tradingSwapContractAddr string, abiFile *abi.ABI, to *artemis_trading_types.TradeOutcome) web3_actions.SendContractTxPayload {
	isToken0 := false
	pairContractAddr, tkn0, _ := artemis_utils.CreateV2TradingPair(to.AmountInAddr, to.AmountOutAddr)
	if tkn0.String() == to.AmountInAddr.String() {
		isToken0 = true
	}
	params := web3_actions.SendContractTxPayload{
		SmartContractAddr: tradingSwapContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       abiFile,
		MethodName:        execSmartContractTradingSwap,
		Params:            []interface{}{pairContractAddr, to.AmountInAddr.String(), isToken0, to.AmountIn.String(), to.AmountOut.String()},
	}
	return params
}

func ExecSmartContractTradingSwap(ctx context.Context, w3c web3_actions.Web3Actions, tradingContractAddr string, abiFile *abi.ABI, to *artemis_trading_types.TradeOutcome) (*types.Transaction, error) {
	scInfo := GetRawDawgSwapAbiPayload(tradingContractAddr, abiFile, to)
	signedTx, err := w3c.GetSignedTxToCallFunctionWithArgs(ctx, &scInfo)
	if err != nil {
		return nil, err
	}
	err = w3c.SendSignedTransaction(ctx, signedTx)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}
