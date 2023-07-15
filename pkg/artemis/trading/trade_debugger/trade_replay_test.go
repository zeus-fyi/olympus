package artemis_trade_debugger

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	artemis_test_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/test_suite/test_cache"
)

/*
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
*/

// 0x80ae3cc1748c10f42e591783001817b8a56b188eb1867282e396a8d99d583d00

func (t *ArtemisTradeDebuggerTestSuite) TestReplayer() {
	txHash := "0x80ae3cc1748c10f42e591783001817b8a56b188eb1867282e396a8d99d583d00"
	err := t.td.Replay(ctx, txHash, true)
	t.NoError(err)
}

func (t *ArtemisTradeDebuggerTestSuite) TestReadTx() {
	artemis_test_cache.LiveTestNetwork.Dial()
	defer artemis_test_cache.LiveTestNetwork.Close()
	frontrun := "0x0213c1ecd07af84469fdb5f790d5639ac93530fdf1311f2a413170e678856a65"
	txHash := "0x035653cdc672256c3ca1da179b9377f59c7290d98b4421d586e227b1a7489a46"
	backrun := "0xd8a03730fcd49362741e15241a43e4336b06bacef70f20b3bc1b4697e493c155"
	tx, _, err := artemis_test_cache.LiveTestNetwork.C.TransactionByHash(ctx, common.HexToHash(frontrun))
	t.NoError(err)
	fmt.Println("tx.GasFeeCap())", tx.GasFeeCap())
	fmt.Println("tx.GasTipCap()", tx.GasTipCap())
	fmt.Println("tx.GasPrice()", tx.GasPrice())
	fmt.Println("tx.Gas()", tx.Gas())

	tx, _, err = artemis_test_cache.LiveTestNetwork.C.TransactionByHash(ctx, common.HexToHash(txHash))
	t.NoError(err)
	fmt.Println("tx.GasFeeCap())", tx.GasFeeCap())
	fmt.Println("tx.GasTipCap()", tx.GasTipCap())
	fmt.Println("tx.GasPrice()", tx.GasPrice())
	fmt.Println("tx.Gas()", tx.Gas())

	tx, _, err = artemis_test_cache.LiveTestNetwork.C.TransactionByHash(ctx, common.HexToHash(backrun))
	t.NoError(err)
	fmt.Println("tx.GasFeeCap())", tx.GasFeeCap())
	fmt.Println("tx.GasTipCap()", tx.GasTipCap())
	fmt.Println("tx.GasPrice()", tx.GasPrice())
	fmt.Println("tx.Gas()", tx.Gas())

	//rx, err := artemis_test_cache.LiveTestNetwork.C.TransactionReceipt(ctx, common.HexToHash(txHash))
	//t.NoError(err)
	//fmt.Println(rx.Status)
	//fmt.Println(rx.GasUsed)
	//fmt.Println(rx.CumulativeGasUsed)
	//fmt.Println(rx.Logs)
}
