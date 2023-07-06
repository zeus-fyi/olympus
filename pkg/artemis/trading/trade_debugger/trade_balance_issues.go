package artemis_trade_debugger

import (
	"fmt"

	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
)

func (t *TradeDebugger) analyzeToken(address accounts.Address) {
	token := address.String()
	den := artemis_trading_cache.TokenMap[token].TransferTaxDenominator
	num := artemis_trading_cache.TokenMap[token].TransferTaxNumerator
	fmt.Println("token: ", token, "tradingTax: ", den, "num: ", num)
}
