package artemis_trading_auxiliary

import "github.com/ethereum/go-ethereum/core/types"

type TxWithMetadata struct {
	TradeType string
	Tx        *types.Transaction
}
