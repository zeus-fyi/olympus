package artemis_trade_debugger

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	artemis_test_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/test_suite/test_cache"
)

func GasAnalysis(ctx context.Context, txHash string) error {
	artemis_test_cache.LiveTestNetwork.Dial()
	defer artemis_test_cache.LiveTestNetwork.Close()
	tx, _, err := artemis_test_cache.LiveTestNetwork.C.TransactionByHash(ctx, common.HexToHash(txHash))
	if err != nil {
		return err
	}
	fmt.Println("tx.Type()", tx.Type())
	fmt.Println("tx.GasFeeCap())", tx.GasFeeCap())
	fmt.Println("tx.GasTipCap()", tx.GasTipCap())
	fmt.Println("tx.GasPrice()", tx.GasPrice())
	fmt.Println("tx.Gas()", tx.Gas())

	//gasFeeAndTip := artemis_eth_units.AddBigInt(tx.GasFeeCap(), tx.GasTipCap())
	//fmt.Println("gasFeeAndTip", gasFeeAndTip)
	rx, err := artemis_test_cache.LiveTestNetwork.C.TransactionReceipt(ctx, common.HexToHash(txHash))
	if err != nil {
		return err
	}
	fmt.Println("rx.GasUsed", rx.GasUsed)
	fmt.Println("rx.EffectiveGasPrice", rx.EffectiveGasPrice.String())
	fmt.Println("rx.TransactionIndex", rx.TransactionIndex)
	//fmt.Println("rx.CumulativeGasUsed", rx.CumulativeGasUsed)

	block, err := artemis_test_cache.LiveTestNetwork.C.BlockByNumber(ctx, rx.BlockNumber)
	if err != nil {
		return err
	}

	fmt.Println("block.BaseFee", block.BaseFee())
	fmt.Println("block.GasLimit", block.GasLimit())
	return nil
}
