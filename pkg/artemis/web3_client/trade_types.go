package web3_client

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
)

type Trade struct {
	TradeMethod                         string `json:"tradeMethod"`
	*JSONSwapETHForExactTokensParams    `json:"swapETHForExactTokensParams,omitempty"`
	*JSONSwapTokensForExactTokensParams `json:"swapTokensForExactTokensParams,omitempty"`
	*JSONSwapExactTokensForTokensParams `json:"swapExactTokensForTokensParams,omitempty"`
	*JSONSwapExactETHForTokensParams    `json:"swapExactETHForTokensParams,omitempty"`
	*JSONSwapExactTokensForETHParams    `json:"swapExactTokensForETHParams,omitempty"`
	*JSONSwapTokensForExactETHParams    `json:"swapTokensForExactETHParams,omitempty"`

	// universal router
	*JSONV2SwapExactInParams  `json:"v2SwapExactInParams,omitempty"`
	*JSONV2SwapExactOutParams `json:"v2SwapExactOutParams,omitempty"`

	*JSONSwapExactTokensForTokensSupportingFeeOnTransferTokensParams `json:"swapExactTokensForTokensSupportingFeeOnTransferTokensParams,omitempty"`
	*JSONSwapExactETHForTokensSupportingFeeOnTransferTokensParams    `json:"swapExactETHForTokensSupportingFeeOnTransferTokensParams,omitempty"`
	*JSONSwapExactTokensForETHSupportingFeeOnTransferTokensParams    `json:"swapExactTokensForETHSupportingFeeOnTransferTokensParams,omitempty"`
	*JSONExactInputParams                                            `json:"exactInputParams,omitempty"`
	*JSONExactOutputParams                                           `json:"exactOutputParams,omitempty"`
	*JSONSwapExactInputSingleArgs                                    `json:"swapExactInputSingleArgs,omitempty"`
	*JSONSwapExactOutputSingleArgs                                   `json:"swapExactOutputSingleArgs,omitempty"`
	*JSONV3SwapExactInParams                                         `json:"v3SwapExactInParams,omitempty"`
	*JSONV3SwapExactOutParams                                        `json:"v3SwapExactOutParams,omitempty"`
	*JSONSwapExactTokensForTokensParamsV3                            `json:"swapExactTokensForTokensParamsV3,omitempty"`
}

type TradeExecutionFlowJSON struct {
	CurrentBlockNumber *big.Int                               `json:"currentBlockNumber"`
	Tx                 artemis_trading_types.JSONTx           `json:"tx"`
	Trade              Trade                                  `json:"trade"`
	InitialPairV3      *uniswap_pricing.JSONUniswapPoolV3     `json:"initialPairV3,omitempty"`
	InitialPair        *uniswap_pricing.JSONUniswapV2Pair     `json:"initialPair,omitempty"`
	FrontRunTrade      artemis_trading_types.JSONTradeOutcome `json:"frontRunTrade"`
	UserTrade          artemis_trading_types.JSONTradeOutcome `json:"userTrade"`
	SandwichTrade      artemis_trading_types.JSONTradeOutcome `json:"sandwichTrade"`
	SandwichPrediction JSONSandwichTradePrediction            `json:"sandwichPrediction"`
}

func (t *TradeExecutionFlowJSON) ConvertToBigIntType() TradeExecutionFlow {
	var p2Pair *uniswap_pricing.UniswapV2Pair
	if t.InitialPair != nil {
		p2Pair = t.InitialPair.ConvertToBigIntType()
	}
	var p3Pair *uniswap_pricing.UniswapV3Pair
	if t.InitialPairV3 != nil {
		p3Pair = t.InitialPairV3.ConvertToBigIntType()
	}

	txConv, err := t.Tx.ConvertToTx()
	if err != nil {
		panic(err)
	}
	return TradeExecutionFlow{
		CurrentBlockNumber: t.CurrentBlockNumber,
		Tx:                 txConv,
		Trade:              t.Trade,
		InitialPair:        p2Pair,
		InitialPairV3:      p3Pair,
		FrontRunTrade:      t.FrontRunTrade.ConvertToBigIntType(),
		UserTrade:          t.UserTrade.ConvertToBigIntType(),
		SandwichTrade:      t.SandwichTrade.ConvertToBigIntType(),
		SandwichPrediction: t.SandwichPrediction.ConvertToBigIntType(),
	}
}

type TradeExecutionFlow struct {
	CurrentBlockNumber *big.Int                           `json:"currentBlockNumber"`
	Tx                 *types.Transaction                 `json:"tx"`
	Trade              Trade                              `json:"trade"`
	InitialPair        *uniswap_pricing.UniswapV2Pair     `json:"initialPair,omitempty"`
	InitialPairV3      *uniswap_pricing.UniswapV3Pair     `json:"initialPairV3,omitempty"`
	FrontRunTrade      artemis_trading_types.TradeOutcome `json:"frontRunTrade"`
	UserTrade          artemis_trading_types.TradeOutcome `json:"userTrade"`
	SandwichTrade      artemis_trading_types.TradeOutcome `json:"sandwichTrade"`
	SandwichPrediction SandwichTradePrediction            `json:"sandwichPrediction"`
}

func (t *TradeExecutionFlow) GetAggregateGasUsage(ctx context.Context, w Web3Client) error {
	//err := t.FrontRunTrade.GetGasUsageForAllTxs(ctx, w)
	//if err != nil {
	//	log.Err(err).Msg("error getting gas usage for front run trade")
	//	return err
	//}
	//err = t.UserTrade.GetGasUsageForAllTxs(ctx, w)
	//if err != nil {
	//	log.Err(err).Msg("error getting gas usage for user trade")
	//	return err
	//}
	//err = t.SandwichTrade.GetGasUsageForAllTxs(ctx, w)
	//if err != nil {
	//	log.Err(err).Msg("error getting gas usage for sandwich trade")
	//	return err
	//}
	return nil
}
