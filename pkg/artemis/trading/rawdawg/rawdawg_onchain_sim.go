package artemis_rawdawg_contract

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi"
	artemis_utils "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/utils"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	web3_actions "github.com/zeus-fyi/zeus/pkg/artemis/web3/client"
)

//type RawDawgV2SimSwapParams struct {
//	Pair      common.Address `json:"_pair"`
//	TokenIn   common.Address `json:"_token_in"`
//	TokenOut  common.Address `json:"_token_out"`
//	AmountIn  *big.Int       `json:"_amountIn"`
//	AmountOut *big.Int       `json:"_amountOut"`
//	IsToken0  bool           `json:"_isToken0"`
//}

const (
	simulateV2AndRevertSwap = "simulateV2AndRevertSwap"
)

func GetRawDawgV2SimSwapAbiPayload(ctx context.Context, tradingSwapContractAddr string, abiFile *abi.ABI, to *artemis_trading_types.TradeOutcome) *web3_actions.SendContractTxPayload {
	isToken0 := false
	pairContractAddr, tkn0, _ := artemis_utils.CreateV2TradingPair(to.AmountInAddr, to.AmountOutAddr)
	if tkn0.Hex() == to.AmountInAddr.Hex() {
		isToken0 = true
	}

	params := &web3_actions.SendContractTxPayload{
		SmartContractAddr: tradingSwapContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       abiFile,
		MethodName:        simulateV2AndRevertSwap,
		Params:            []interface{}{pairContractAddr, to.AmountInAddr, to.AmountOutAddr, isToken0, to.AmountIn, to.AmountOut},
	}
	return params
}
