package artemis_trade_debugger

import (
	"errors"
	"fmt"

	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_realtime_trading "github.com/zeus-fyi/olympus/pkg/artemis/trading"
	"github.com/zeus-fyi/olympus/pkg/artemis/trading/async_analysis"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

type TradeDebugger struct {
	insertNewTxs      bool
	onlyWETH          bool
	dat               artemis_realtime_trading.ActiveTrading
	ContractAnalysis  async_analysis.ContractAnalysis
	LiveNetworkClient web3_client.Web3Client
}

func NewTradeDebugger(a artemis_realtime_trading.ActiveTrading, lnc web3_client.Web3Client) TradeDebugger {
	return TradeDebugger{
		dat:               a,
		LiveNetworkClient: lnc,
	}
}

func NewTradeDebuggerWorkflowAnalysis(a artemis_realtime_trading.ActiveTrading, lnc web3_client.Web3Client) TradeDebugger {
	td := TradeDebugger{
		dat:               a,
		LiveNetworkClient: lnc,
	}
	td.insertNewTxs = true
	return td
}

type HistoricalAnalysisDebug struct {
	HistoricalAnalysis artemis_mev_models.HistoricalAnalysis
	TradePrediction    web3_client.TradeExecutionFlow
	TradeParams        any
}

/*
type HistoricalAnalysis struct {
	artemis_autogen_bases.EthMempoolMevTx
	artemis_autogen_bases.EthMevTxAnalysis
}
type EthMempoolMevTx struct {
	TxID              int    `db:"tx_id" json:"txID"`
	To                string `db:"to" json:"to"`
	ProtocolNetworkID int    `db:"protocol_network_id" json:"protocolNetworkID"`
	TxFlowPrediction  string `db:"tx_flow_prediction" json:"txFlowPrediction"`
	TxHash            string `db:"tx_hash" json:"txHash"`
	Nonce             int    `db:"nonce" json:"nonce"`
	From              string `db:"from" json:"from"`
	BlockNumber       int    `db:"block_number" json:"blockNumber"`
	Tx                string `db:"tx" json:"tx"`
}
*/

func (h *HistoricalAnalysisDebug) GetBlockNumber() int {
	return h.HistoricalAnalysis.BlockNumber
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
