package web3_client

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
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

	*JSONV3SwapExactInParams  `json:"v3SwapExactInParams,omitempty"`
	*JSONV3SwapExactOutParams `json:"v3SwapExactOutParams,omitempty"`
}

type TradeExecutionFlowJSON struct {
	CurrentBlockNumber *big.Int                    `json:"currentBlockNumber"`
	Tx                 *types.Transaction          `json:"tx"`
	Trade              Trade                       `json:"trade"`
	InitialPair        JSONUniswapV2Pair           `json:"initialPair"`
	FrontRunTrade      JSONTradeOutcome            `json:"frontRunTrade"`
	UserTrade          JSONTradeOutcome            `json:"userTrade"`
	SandwichTrade      JSONTradeOutcome            `json:"sandwichTrade"`
	SandwichPrediction JSONSandwichTradePrediction `json:"sandwichPrediction"`
}

func (t *TradeExecutionFlowJSON) ConvertToBigIntType() TradeExecutionFlow {
	return TradeExecutionFlow{
		CurrentBlockNumber: t.CurrentBlockNumber,
		Tx:                 t.Tx,
		Trade:              t.Trade,
		InitialPair:        t.InitialPair.ConvertToBigIntType(),
		FrontRunTrade:      t.FrontRunTrade.ConvertToBigIntType(),
		UserTrade:          t.UserTrade.ConvertToBigIntType(),
		SandwichTrade:      t.SandwichTrade.ConvertToBigIntType(),
		SandwichPrediction: t.SandwichPrediction.ConvertToBigIntType(),
	}
}

type TradeExecutionFlow struct {
	CurrentBlockNumber *big.Int                `json:"currentBlockNumber"`
	Tx                 *types.Transaction      `json:"tx"`
	Trade              Trade                   `json:"trade"`
	InitialPair        UniswapV2Pair           `json:"initialPair"`
	FrontRunTrade      TradeOutcome            `json:"frontRunTrade"`
	UserTrade          TradeOutcome            `json:"userTrade"`
	SandwichTrade      TradeOutcome            `json:"sandwichTrade"`
	SandwichPrediction SandwichTradePrediction `json:"sandwichPrediction"`
}

func (t *TradeExecutionFlow) GetAggregateGasUsage(ctx context.Context, w Web3Client) error {
	err := t.FrontRunTrade.GetGasUsageForAllTxs(ctx, w)
	if err != nil {
		log.Err(err).Msg("error getting gas usage for front run trade")
		return err
	}
	err = t.UserTrade.GetGasUsageForAllTxs(ctx, w)
	if err != nil {
		log.Err(err).Msg("error getting gas usage for user trade")
		return err
	}
	err = t.SandwichTrade.GetGasUsageForAllTxs(ctx, w)
	if err != nil {
		log.Err(err).Msg("error getting gas usage for sandwich trade")
		return err
	}
	return nil
}
