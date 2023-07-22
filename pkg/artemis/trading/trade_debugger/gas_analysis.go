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

	rx, err := artemis_test_cache.LiveTestNetwork.C.TransactionReceipt(ctx, common.HexToHash(txHash))
	if err != nil {
		return err
	}
	fmt.Println("gasUsed", rx.GasUsed)
	fmt.Println("rx.CumulativeGasUsed", rx.CumulativeGasUsed)
	return nil
}
