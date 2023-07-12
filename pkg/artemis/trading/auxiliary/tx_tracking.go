package artemis_trading_auxiliary

import (
	"fmt"
)

func (a *AuxiliaryTradingUtils) trackTxs(txs MevTxGroup) {
	for _, tx := range txs.OrderedTxs {
		fmt.Println("tx", tx.Tx.Hash().Hex())
	}
}
