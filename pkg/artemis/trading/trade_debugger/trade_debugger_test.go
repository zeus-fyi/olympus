package artemis_trade_debugger

import (
	"context"
	"fmt"
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
	td := NewTradeDebugger(at, &uni, t.MainnetWeb3User)
	t.Require().NotEmpty(td)
	t.td = td
}

func (t *ArtemisTradeDebuggerTestSuite) TestReplayDebugger1() {
	txHash := "0x5327295e1ed6d59faaf98d04697b0316fb8ad4b767d2e7f5addb3981c3b5d3b7"
	hist, err := t.td.lookupMevTx(ctx, txHash)
	t.Require().Nil(err)
	t.Require().NotEmpty(hist)
	for _, h := range hist {
		fmt.Println(h.HistoricalAnalysis.TradeMethod)
		fmt.Println(h.HistoricalAnalysis.EndReason)
		tf, serr := h.BinarySearch()
		t.Require().Nil(serr)
		t.Require().NotEmpty(tf)
	}
	t.Assert().NotEmpty(hist)
}

func (t *ArtemisTradeDebuggerTestSuite) TestReplayDebugger2() {
	txHash := "0x5327295e1ed6d59faaf98d04697b0316fb8ad4b767d2e7f5addb3981c3b5d3b7"
	hist, err := t.td.lookupMevTx(ctx, txHash)
	t.Require().Nil(err)
	t.Require().NotEmpty(hist)
	for _, h := range hist {
		err = t.td.ResetNetwork(&h.TradePrediction)
		t.Require().Nil(err)
	}
	t.Assert().NotEmpty(hist)
}

func TestArtemisTradeDebuggerTestSuite(t *testing.T) {
	suite.Run(t, new(ArtemisTradeDebuggerTestSuite))
}
