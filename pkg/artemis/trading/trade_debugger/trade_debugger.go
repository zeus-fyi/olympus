package artemis_trade_debugger

import (
	"errors"
	"fmt"

	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_realtime_trading "github.com/zeus-fyi/olympus/pkg/artemis/trading"
	"github.com/zeus-fyi/olympus/pkg/artemis/trading/async_analysis"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
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
func BinarySearch(tf web3_client.TradeExecutionFlow) (web3_client.TradeExecutionFlow, error) {
	switch tf.Trade.TradeMethod {
	case artemis_trading_constants.V2SwapExactIn:
		params := tf.Trade.JSONV2SwapExactInParams
		paramsBigInt, err := params.ConvertToBigIntType()
		if err != nil {
			return tf, fmt.Errorf("error in binary search: %w", err)
		}
		if tf.InitialPair == nil {
			return tf, errors.New("initial pair is nil")
		}
		search, err := paramsBigInt.BinarySearch(*tf.InitialPair)
		if err != nil {
			return tf, fmt.Errorf("error in binary search: %w", err)
		}
		return search, nil
	case artemis_trading_constants.SwapExactTokensForTokensSupportingFeeOnTransferTokens:
		//params := tf.Trade.JSONSwapExactTokensForTokensSupportingFeeOnTransferTokensParams
		//vals, err := params.ConvertToBigIntType()
		//if err != nil {
		//	return tf, fmt.Errorf("error in binary search: %w", err)
		//}
		//if tf.InitialPair == nil {
		//	return tf, errors.New("initial pair is nil")
		//}
		//search, err := vals.BinarySearch(*tf.InitialPair)
		//if err != nil {
		//	return tf, fmt.Errorf("error in binary search: %w", err)
		//}
		//return search, nil
	case artemis_trading_constants.SwapExactETHForTokensSupportingFeeOnTransferTokens:
		//params := tf.Trade.JSONSwapExactETHForTokensSupportingFeeOnTransferTokensParams
		//if tf.InitialPair == nil {
		//	return tf, errors.New("initial pair is nil")
		//}
		//if tf.InitialPair == nil {
		//	return tf, errors.New("initial pair is nil")
		//}
		//search, err := params.BinarySearch(*tf.InitialPair)
		//if err != nil {
		//	return tf, fmt.Errorf("error in binary search: %w", err)
		//}
		//return search, nil
	case artemis_trading_constants.SwapExactTokensForETHSupportingFeeOnTransferTokens:
		//params := tf.Trade.JSONSwapExactTokensForETHSupportingFeeOnTransferTokensParams
		//if tf.InitialPair == nil {
		//	return tf, errors.New("initial pair is nil")
		//}
		//if tf.InitialPair == nil {
		//	return tf, errors.New("initial pair is nil")
		//}
		//search, err := params.BinarySearch(*tf.InitialPair)
		//if err != nil {
		//	return tf, fmt.Errorf("error in binary search: %w", err)
		//}
		//return search, nil
	case artemis_trading_constants.SwapExactTokensForTokens:
		//params := tf.Trade.JSONSwapExactTokensForTokensParams
		//if tf.InitialPair == nil {
		//	return tf, errors.New("initial pair is nil")
		//}
		//if tf.InitialPair == nil {
		//	return tf, errors.New("initial pair is nil")
		//}
		//search, err := params.BinarySearch(*tf.InitialPair)
		//if err != nil {
		//	return tf, fmt.Errorf("error in binary search: %w", err)
		//}
		//return search, nil
	case artemis_trading_constants.SwapExactETHForTokens:
		params := tf.Trade.JSONSwapExactETHForTokensParams
		if tf.InitialPair == nil {
			return tf, errors.New("initial pair is nil")
		}
		paramsBigInt := params.ConvertToBigIntType()
		search, err := paramsBigInt.BinarySearch(*tf.InitialPair)
		if err != nil {
			return tf, fmt.Errorf("error in binary search: %w", err)
		}
		return search, nil
	case artemis_trading_constants.SwapExactTokensForETH:
		//params := tf.Trade.JSONSwapExactTokensForETHParams
		//if tf.InitialPair == nil {
		//	return tf, errors.New("initial pair is nil")
		//}
		//search, err := params.BinarySearch(*tf.InitialPair)
		//if err != nil {
		//	return tf, fmt.Errorf("error in binary search: %w", err)
		//}
		//return search, nil
	case artemis_trading_constants.SwapTokensForExactTokens:
		//params := tf.Trade.JSONSwapTokensForExactTokensParams
		//if tf.InitialPair == nil {
		//	return tf, errors.New("initial pair is nil")
		//}
		//search, err := params.BinarySearch(*tf.InitialPair)
		//if err != nil {
		//	return tf, fmt.Errorf("error in binary search: %w", err)
		//}
		//return search, nil
	case artemis_trading_constants.SwapTokensForExactETH:
		//params := tf.Trade.JSONSwapTokensForExactETHParams
		//if tf.InitialPair == nil {
		//	return tf, errors.New("initial pair is nil")
		//}
		//search, err := params.BinarySearch(*tf.InitialPair)
		//if err != nil {
		//	return tf, fmt.Errorf("error in binary search: %w", err)
		//}
		//return search, nil
	case artemis_trading_constants.SwapETHForExactTokens:
		params := tf.Trade.JSONSwapETHForExactTokensParams
		if tf.InitialPair == nil {
			return tf, errors.New("initial pair is nil")
		}
		paramsBigInt := params.ConvertToBigIntType()
		search, err := paramsBigInt.BinarySearch(*tf.InitialPair)
		if err != nil {
			return tf, fmt.Errorf("error in binary search: %w", err)
		}
		return search, nil
	default:
		fmt.Println(tf.Trade.TradeMethod, "tradeMethod not supported for binary search historical analysis")
	}
	return web3_client.TradeExecutionFlow{}, errors.New("no trade params found")
}
