package artemis_trading_auxiliary

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
)

func (a *AuxiliaryTradingUtils) trackTxs(txs ...*types.Transaction) {
	for _, tx := range txs {
		fmt.Println("tx", tx.Hash().String())
	}
}
