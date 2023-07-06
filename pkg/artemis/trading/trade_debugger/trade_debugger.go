package artemis_trade_debugger

import (
	"errors"
	"fmt"
	"math/big"

	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_realtime_trading "github.com/zeus-fyi/olympus/pkg/artemis/trading"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

type TradeDebugger struct {
	UniswapClient *web3_client.UniswapClient
	ActiveTrading artemis_realtime_trading.ActiveTrading
}

func NewTradeDebugger(a artemis_realtime_trading.ActiveTrading, u *web3_client.UniswapClient) TradeDebugger {
	return TradeDebugger{
		ActiveTrading: a,
		UniswapClient: u,
	}
}

type HistoricalAnalysisDebug struct {
	HistoricalAnalysis artemis_mev_models.HistoricalAnalysis
	TradePrediction    web3_client.TradeExecutionFlow
	TradeParams        any
}

func (h *HistoricalAnalysisDebug) BinarySearch() (web3_client.TradeExecutionFlow, error) {
	tf := web3_client.TradeExecutionFlow{}
	v := h.TradeParams
	switch v.(type) {
	case *web3_client.V2SwapExactInParams:
		params := v.(*web3_client.V2SwapExactInParams)
		if h.TradePrediction.InitialPair == nil {
			return tf, errors.New("initial pair is nil")
		}
		search := params.BinarySearch(*h.TradePrediction.InitialPair)
		return search.ConvertToBigIntTypeWithoutTx(), nil
	case *web3_client.SwapExactTokensForTokensSupportingFeeOnTransferTokensParams:
		params := v.(*web3_client.SwapExactTokensForTokensSupportingFeeOnTransferTokensParams)
		if h.TradePrediction.InitialPair == nil {
			return tf, errors.New("initial pair is nil")
		}
		search := params.BinarySearch(*h.TradePrediction.InitialPair)
		return search.ConvertToBigIntTypeWithoutTx(), nil
	case *web3_client.SwapExactETHForTokensSupportingFeeOnTransferTokensParams:
		params := v.(*web3_client.SwapExactETHForTokensSupportingFeeOnTransferTokensParams)
		if h.TradePrediction.InitialPair == nil {
			return tf, errors.New("initial pair is nil")
		}
		search := params.BinarySearch(*h.TradePrediction.InitialPair)
		return search.ConvertToBigIntTypeWithoutTx(), nil
	case *web3_client.SwapExactTokensForETHSupportingFeeOnTransferTokensParams:
		params := v.(*web3_client.SwapExactTokensForETHSupportingFeeOnTransferTokensParams)
		if h.TradePrediction.InitialPair == nil {
			return tf, errors.New("initial pair is nil")
		}
		search := params.BinarySearch(*h.TradePrediction.InitialPair)
		return search.ConvertToBigIntTypeWithoutTx(), nil
	case *web3_client.SwapExactTokensForTokensParams:
		params := v.(*web3_client.SwapExactTokensForTokensParams)
		if h.TradePrediction.InitialPair == nil {
			return tf, errors.New("initial pair is nil")
		}
		search := params.BinarySearch(*h.TradePrediction.InitialPair)
		return search.ConvertToBigIntTypeWithoutTx(), nil
	case *web3_client.SwapExactETHForTokensParams:
		params := v.(*web3_client.SwapExactETHForTokensParams)
		if h.TradePrediction.InitialPair == nil {
			return tf, errors.New("initial pair is nil")

		}
		search := params.BinarySearch(*h.TradePrediction.InitialPair)
		return search.ConvertToBigIntTypeWithoutTx(), nil
	case *web3_client.SwapExactTokensForETHParams:
		params := v.(*web3_client.SwapExactTokensForETHParams)
		if h.TradePrediction.InitialPair == nil {
			return tf, errors.New("initial pair is nil")

		}
		search := params.BinarySearch(*h.TradePrediction.InitialPair)
		return search.ConvertToBigIntTypeWithoutTx(), nil
	case *web3_client.SwapTokensForExactTokensParams:
		params := v.(*web3_client.SwapTokensForExactTokensParams)
		if h.TradePrediction.InitialPair == nil {
			return tf, errors.New("initial pair is nil")
		}
		search := params.BinarySearch(*h.TradePrediction.InitialPair)
		return search.ConvertToBigIntTypeWithoutTx(), nil
	case *web3_client.SwapTokensForExactETHParams:
		params := v.(*web3_client.SwapTokensForExactETHParams)
		if h.TradePrediction.InitialPair == nil {
			return tf, errors.New("initial pair is nil")
		}
		search := params.BinarySearch(*h.TradePrediction.InitialPair)
		return search.ConvertToBigIntTypeWithoutTx(), nil
	case *web3_client.SwapETHForExactTokensParams:
		params := v.(*web3_client.SwapETHForExactTokensParams)
		if h.TradePrediction.InitialPair == nil {
			return tf, errors.New("initial pair is nil")
		}
		search := params.BinarySearch(*h.TradePrediction.InitialPair)
		return search.ConvertToBigIntTypeWithoutTx(), nil
	default:
		fmt.Println(v)
	}
	return web3_client.TradeExecutionFlow{}, errors.New("no trade params found")
}

/*
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
*/

func ParseBigInt(i interface{}) (*big.Int, error) {
	switch v := i.(type) {
	case *big.Int:
		return i.(*big.Int), nil
	case string:
		base := 10
		result := new(big.Int)
		_, ok := result.SetString(v, base)
		if !ok {
			return nil, fmt.Errorf("failed to parse string '%s' into big.Int", v)
		}
		return result, nil
	case uint32:
		return big.NewInt(int64(v)), nil
	case int64:
		return big.NewInt(v), nil
	default:
		return nil, fmt.Errorf("input is not a string or int64")
	}
}
