package artemis_rawdawg_contract

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_v3/constants"
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

type QuoteRawDawgSwapParams struct {
	TokenIn           accounts.Address `json:"tokenIn"`
	TokenOut          accounts.Address `json:"tokenOut"`
	Fee               *big.Int         `json:"fee"`
	AmountIn          *big.Int         `json:"amountIn"`
	SqrtPriceLimitX96 *big.Int         `json:"sqrtPriceLimitX96"`
}

func GetRawDawgV2SimSwapAbiPayload(ctx context.Context, tradingSwapContractAddr string, abiFile *abi.ABI, to *artemis_trading_types.TradeOutcome) *web3_actions.SendContractTxPayload {
	//isToken0 := false
	//pairContractAddr, tkn0, _ := artemis_utils.CreateV2TradingPair(to.AmountInAddr, to.AmountOutAddr)
	//if tkn0.Hex() == to.AmountInAddr.Hex() {
	//	isToken0 = true
	//}
	//
	//tfp := artemis_trading_types.TokenFeePath{
	//	TokenIn: to.AmountInAddr,
	//	Path: []artemis_trading_types.TokenFee{
	//		{
	//			Token: to.AmountOutAddr,
	//			Fee:   artemis_eth_units.NewBigInt(int(constants.FeeMedium)),
	//		},
	//	},
	//}

	qp := QuoteRawDawgSwapParams{
		TokenIn:           to.AmountInAddr,
		TokenOut:          to.AmountOutAddr,
		Fee:               artemis_eth_units.NewBigInt(int(constants.FeeMedium)),
		AmountIn:          to.AmountIn,
		SqrtPriceLimitX96: artemis_eth_units.NewBigInt(0),
	}
	params := &web3_actions.SendContractTxPayload{
		SmartContractAddr: tradingSwapContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       abiFile,
		MethodName:        simulateV2AndRevertSwap,
		Params:            []interface{}{qp},
	}
	return params
}
