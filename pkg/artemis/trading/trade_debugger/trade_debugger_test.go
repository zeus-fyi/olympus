package artemis_trade_debugger

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	artemis_realtime_trading "github.com/zeus-fyi/olympus/pkg/artemis/trading"
	artemis_trading_test_suite "github.com/zeus-fyi/olympus/pkg/artemis/trading/test_suite"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

type ArtemisTradeDebuggerTestSuite struct {
	artemis_trading_test_suite.ArtemisTradingTestSuite
}

var ctx = context.Background()

func (t *ArtemisTradeDebuggerTestSuite) TestDebugger() {
	uni := web3_client.InitUniswapClient(ctx, t.ProxiedMainnetUser)
	at := artemis_realtime_trading.NewActiveTradingModuleWithoutMetrics(&uni)
	td := NewTradeDebugger(at, &uni)
	t.Require().NotEmpty(td)

	txHash := "0x5327295e1ed6d59faaf98d04697b0316fb8ad4b767d2e7f5addb3981c3b5d3b7"
	err := td.GetTxFromHash(ctx, txHash)
	t.Require().Nil(err)

}

func TestArtemisTradeDebuggerTestSuite(t *testing.T) {
	suite.Run(t, new(ArtemisTradeDebuggerTestSuite))
}
