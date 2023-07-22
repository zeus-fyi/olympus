package artemis_trade_debugger

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	artemis_trading_auxiliary "github.com/zeus-fyi/olympus/pkg/artemis/trading/auxiliary"
	artemis_test_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/test_suite/test_cache"
)

func (t *ArtemisTradeDebuggerTestSuite) TestGasAnalysis() {
	txHash := "0xf2545a49a3e3e3e8fa0ec699943b1f454adda02888a16c39a64774bbdab248eb"
	err := GasAnalysis(ctx, txHash)
	t.NoError(err)

	//tx, _, err := artemis_test_cache.LiveTestNetwork.C.TransactionByHash(ctx, common.HexToHash(frontrun))
	artemis_trading_auxiliary.ApplyFrontRunGasAdjustment(nil)
	artemis_trading_auxiliary.ApplyBackrunGasAdjustment(nil)
	artemis_trading_auxiliary.ApplyTxType2UserGasAdjustment(nil)
}

func (t *ArtemisTradeDebuggerTestSuite) TestReadTx() {
	artemis_test_cache.LiveTestNetwork.Dial()
	defer artemis_test_cache.LiveTestNetwork.Close()
	frontrun := "0x0213c1ecd07af84469fdb5f790d5639ac93530fdf1311f2a413170e678856a65"
	txHash := "0x035653cdc672256c3ca1da179b9377f59c7290d98b4421d586e227b1a7489a46"
	backrun := "0xd8a03730fcd49362741e15241a43e4336b06bacef70f20b3bc1b4697e493c155"
	tx, _, err := artemis_test_cache.LiveTestNetwork.C.TransactionByHash(ctx, common.HexToHash(frontrun))
	t.NoError(err)
	fmt.Println("tx.Type()", tx.Type())
	fmt.Println("tx.GasFeeCap())", tx.GasFeeCap())
	fmt.Println("tx.GasTipCap()", tx.GasTipCap())
	fmt.Println("tx.GasPrice()", tx.GasPrice())
	fmt.Println("tx.Gas()", tx.Gas())

	rx, err := artemis_test_cache.LiveTestNetwork.C.TransactionReceipt(ctx, common.HexToHash(txHash))
	t.NoError(err)
	fmt.Println("gasUsed", rx.GasUsed)
	fmt.Println("rx.CumulativeGasUsed", rx.CumulativeGasUsed)

	tx, _, err = artemis_test_cache.LiveTestNetwork.C.TransactionByHash(ctx, common.HexToHash(txHash))
	t.NoError(err)
	fmt.Println("tx.Type()", tx.Type())
	fmt.Println("tx.GasFeeCap())", tx.GasFeeCap())
	fmt.Println("tx.GasTipCap()", tx.GasTipCap())
	fmt.Println("tx.GasPrice()", tx.GasPrice())
	fmt.Println("tx.Gas()", tx.Gas())

	rx, err = artemis_test_cache.LiveTestNetwork.C.TransactionReceipt(ctx, common.HexToHash(txHash))
	t.NoError(err)
	fmt.Println("gasUsed", rx.GasUsed)
	fmt.Println("rx.CumulativeGasUsed", rx.CumulativeGasUsed)

	tx, _, err = artemis_test_cache.LiveTestNetwork.C.TransactionByHash(ctx, common.HexToHash(backrun))
	t.NoError(err)
	fmt.Println("tx.Type()", tx.Type())
	fmt.Println("tx.GasFeeCap())", tx.GasFeeCap())
	fmt.Println("tx.GasTipCap()", tx.GasTipCap())
	fmt.Println("tx.GasPrice()", tx.GasPrice())
	fmt.Println("tx.Gas()", tx.Gas())

	rx, err = artemis_test_cache.LiveTestNetwork.C.TransactionReceipt(ctx, common.HexToHash(txHash))
	t.NoError(err)
	fmt.Println("gasUsed", rx.GasUsed)
	fmt.Println("rx.CumulativeGasUsed", rx.CumulativeGasUsed)
}
