package artemis_trade_debugger

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_realtime_trading "github.com/zeus-fyi/olympus/pkg/artemis/trading"
	artemis_trading_test_suite "github.com/zeus-fyi/olympus/pkg/artemis/trading/test_suite"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

type ArtemisTradeDebuggerTestSuite struct {
	artemis_trading_test_suite.ArtemisTradingTestSuite
	td TradeDebugger
}

var ctx = context.Background()

func (t *ArtemisTradeDebuggerTestSuite) SetupTest() {
	t.ArtemisTradingTestSuite.SetupTest()
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	uni := web3_client.InitUniswapClient(ctx, t.ProxiedMainnetUser)
	at := artemis_realtime_trading.NewActiveTradingModuleWithoutMetrics(&uni)
	td := NewTradeDebugger(at, &uni)
	t.Require().NotEmpty(td)
	t.td = td
}
func TestArtemisTradeDebuggerTestSuite(t *testing.T) {
	suite.Run(t, new(ArtemisTradeDebuggerTestSuite))
}
