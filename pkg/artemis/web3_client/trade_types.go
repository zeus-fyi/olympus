package web3_client

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	artemis_flashbots "github.com/zeus-fyi/olympus/pkg/artemis/trading/flashbots"
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

func (t *TradeExecutionFlowJSON) ConvertToBigIntTypeWithoutTx() (TradeExecutionFlow, error) {
	var p2Pair *uniswap_pricing.UniswapV2Pair
	if t.InitialPair != nil {
		p2Pair = t.InitialPair.ConvertToBigIntType()
	}
	var p3Pair *uniswap_pricing.UniswapV3Pair
	if t.InitialPairV3 != nil {
		p3Pair = t.InitialPairV3.ConvertToBigIntType()
	}
	if p2Pair == nil && p3Pair == nil {
		log.Error().Msg("TradeExecutionFlowJSON: failed to convert pair")
		return TradeExecutionFlow{}, errors.New("both pricing pairs are nil")
	}
	sp, err := t.SandwichPrediction.ConvertToBigIntType()
	if err != nil {
		log.Error().Msg("TradeExecutionFlowJSON: failed to convert sandwich prediction")
		return TradeExecutionFlow{}, err
	}
	return TradeExecutionFlow{
		CurrentBlockNumber: t.CurrentBlockNumber,
		Trade:              t.Trade,
		InitialPair:        p2Pair,
		InitialPairV3:      p3Pair,
		FrontRunTrade:      t.FrontRunTrade.ConvertToBigIntType(),
		UserTrade:          t.UserTrade.ConvertToBigIntType(),
		SandwichTrade:      t.SandwichTrade.ConvertToBigIntType(),
		SandwichPrediction: sp,
	}, err
}

func (t *TradeExecutionFlowJSON) ConvertToBigIntType() (TradeExecutionFlow, error) {
	var p2Pair *uniswap_pricing.UniswapV2Pair
	if t.InitialPair != nil {
		p2Pair = t.InitialPair.ConvertToBigIntType()
	}
	var p3Pair *uniswap_pricing.UniswapV3Pair
	if t.InitialPairV3 != nil {
		p3Pair = t.InitialPairV3.ConvertToBigIntType()
	}
	if p2Pair == nil && p3Pair == nil {
		log.Error().Msg("TradeExecutionFlowJSON: failed to convert pair")
		return TradeExecutionFlow{}, errors.New("both pricing pairs are nil")
	}
	txConv, err := t.Tx.ConvertToTx()
	if err != nil {
		log.Error().Msg("TradeExecutionFlowJSON: failed to convert tx")
		return TradeExecutionFlow{}, err
	}
	sp, err := t.SandwichPrediction.ConvertToBigIntType()
	if err != nil {
		log.Error().Msg("TradeExecutionFlowJSON: failed to convert sandwich prediction")
		return TradeExecutionFlow{}, err
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
		SandwichPrediction: sp,
	}, nil
}

func (t *TradeExecutionFlow) ConvertToJSONType() (TradeExecutionFlowJSON, error) {
	newJsonTx := artemis_trading_types.JSONTx{}
	err := newJsonTx.UnmarshalTx(t.Tx)
	if err != nil {
		return TradeExecutionFlowJSON{}, err
	}
	var v3Pair *uniswap_pricing.JSONUniswapPoolV3
	if t.InitialPairV3 != nil {
		v3Pair = t.InitialPairV3.ConvertToJSONType()
	}
	var v2Pair *uniswap_pricing.JSONUniswapV2Pair
	if t.InitialPair != nil {
		v2Pair = t.InitialPair.ConvertToJSONType()
	}
	if v2Pair == nil && v3Pair == nil {
		log.Error().Msg("TradeExecutionFlowJSON: failed to convert pair")
		return TradeExecutionFlowJSON{}, errors.New("both pricing pairs are nil")
	}
	return TradeExecutionFlowJSON{
		CurrentBlockNumber: t.CurrentBlockNumber,
		Tx:                 newJsonTx,
		Trade:              t.Trade,
		InitialPairV3:      v3Pair,
		InitialPair:        v2Pair,
		FrontRunTrade:      t.FrontRunTrade.ConvertToJSONType(),
		UserTrade:          t.UserTrade.ConvertToJSONType(),
		SandwichTrade:      t.SandwichTrade.ConvertToJSONType(),
		SandwichPrediction: t.SandwichPrediction.ConvertToJSONType(),
	}, nil
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
	Bundle             *artemis_flashbots.MevTxBundle     `json:"bundle,omitempty"`
}

func (t *TradeExecutionFlow) AreAllTradesValid() bool {
	if t.CurrentBlockNumber == nil {
		log.Warn().Msg("TradeExecutionFlow: current block number is nil")
		return false
	}
	if t.Tx == nil {
		log.Warn().Msg("TradeExecutionFlow: tx is nil")
		return false
	}
	if t.Tx.To() == nil {
		log.Warn().Msg("TradeExecutionFlow: tx to is nil")
		return false
	}
	if !t.FrontRunTrade.AreTradeParamsValid() {
		log.Warn().Msg("TradeExecutionFlow: front run trade is not valid")
		return false
	}
	if !t.UserTrade.AreTradeParamsValid() {
		log.Warn().Msg("TradeExecutionFlow: user trade is not valid")
		return false
	}
	if !t.SandwichTrade.AreTradeParamsValid() {
		log.Warn().Msg("TradeExecutionFlow: sandwich trade is not valid")
		return false
	}
	if !t.SandwichPrediction.CheckForValidityAndProfit() {
		log.Warn().Msg("TradeExecutionFlow: sandwich prediction is not valid")
		return false
	}
	return true
}

func (t *TradeExecutionFlow) GetAggregateGasUsage(ctx context.Context, w Web3Client) error {
	err := t.FrontRunTrade.GetGasUsageForAllTxs(ctx, w.Web3Actions)
	if err != nil {
		log.Err(err).Msg("error getting gas usage for front run trade")
		return err
	}
	err = t.UserTrade.GetGasUsageForAllTxs(ctx, w.Web3Actions)
	if err != nil {
		log.Err(err).Msg("error getting gas usage for user trade")
		return err
	}
	err = t.SandwichTrade.GetGasUsageForAllTxs(ctx, w.Web3Actions)
	if err != nil {
		log.Err(err).Msg("error getting gas usage for sandwich trade")
		return err
	}
	return nil
}

func UnmarshalTradeExecutionFlow(tfStr string) (TradeExecutionFlowJSON, error) {
	tf := TradeExecutionFlowJSON{}
	by := []byte(tfStr)
	berr := json.Unmarshal(by, &tf)
	if berr != nil {
		log.Err(berr).Msg("error unmarshalling trade execution flow")
		return tf, berr
	}
	return tf, nil
}
