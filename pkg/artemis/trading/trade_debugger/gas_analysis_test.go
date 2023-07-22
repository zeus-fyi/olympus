package artemis_trade_debugger

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	artemis_test_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/test_suite/test_cache"
)

func (t *ArtemisTradeDebuggerTestSuite) TestGasAnalysis() {
	var err error
	jaredFrontRunTx := "0xb50378511787f9d867f10f1bef1bad9216d01883ac2c2fec804a42781061ef7f"
	err = GasAnalysis(ctx, jaredFrontRunTx)
	t.NoError(err)

	//block.BaseFee 13848216283
	//block.GasLimit 30000000
	/*
		zeus
			gas tip cap 0
			gas fee cap 27746432566
			gas limit 231372
		vs

		jared
			gas tip cap 13848216283
			gas fee cap 13848216283
			gas limit 555787
	*/

	fmt.Println("====================================================================")
	//txHash := "0xf2545a49a3e3e3e8fa0ec699943b1f454adda02888a16c39a64774bbdab248eb"
	//err = GasAnalysis(ctx, txHash)
	//t.NoError(err)
	// 0x0102d8ba5d8e031eabc28b17f1b94a8440d36baad3c6381fa21907854d487fe5
	txHash := "0x0102d8ba5d8e031eabc28b17f1b94a8440d36baad3c6381fa21907854d487fe5"
	err = GasAnalysis(ctx, txHash)
	t.NoError(err)

	fmt.Println("====================================================================")
	backRun := "0x1686c5cbaebde990d638dcd681f9a5a4814ed5025129c5dc1d6c2b94a4177388"

	err = GasAnalysis(ctx, backRun)
	t.NoError(err)

	/*
		zeus
			gas tip cap 110985730264
			gas fee cap 110985730264
			gas limit 462744
		vs

		jared
			gas tip cap 13848216283
			gas fee cap 13848216283
			gas limit 555787
	*/

	//tx, _, err := artemis_test_cache.LiveTestNetwork.C.TransactionByHash(ctx, common.HexToHash(frontrun))
	//artemis_trading_auxiliary.ApplyFrontRunGasAdjustment(nil)
	//artemis_trading_auxiliary.ApplyBackrunGasAdjustment(nil)
	//artemis_trading_auxiliary.ApplyTxType2UserGasAdjustment(nil)
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
