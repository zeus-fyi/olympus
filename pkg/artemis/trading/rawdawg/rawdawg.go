package artemis_rawdawg_contract

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
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

func GetRawdawgSwapAbiPayload(tradingSwapContractAddr, pairContractAddr string, to *artemis_trading_types.TradeOutcome, isToken0 bool) web3_actions.SendContractTxPayload {
	params := web3_actions.SendContractTxPayload{
		SmartContractAddr: tradingSwapContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       RawdawgAbi,
		MethodName:        execSmartContractTradingSwap,
		Params:            []interface{}{pairContractAddr, to.AmountInAddr.String(), isToken0, to.AmountIn.String(), to.AmountOut.String()},
	}
	return params
}

func ExecSmartContractTradingSwap(ctx context.Context, w3c web3_actions.Web3Actions, tradingContractAddr string, pair uniswap_pricing.UniswapV2Pair, to *artemis_trading_types.TradeOutcome) (*types.Transaction, error) {
	tokenNum := pair.GetTokenNumber(to.AmountInAddr)
	scInfo := GetRawdawgSwapAbiPayload(tradingContractAddr, pair.PairContractAddr, to, tokenNum == 0)
	// TODO implement better gas estimation
	scInfo.GasLimit = 3000000
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
