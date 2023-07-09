package artemis_trading_auxiliary

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
)

func (a *AuxiliaryTradingUtils) trackTx(tx *types.Transaction) {
	// todo
	fmt.Println("tx", tx.Hash().String())
}
