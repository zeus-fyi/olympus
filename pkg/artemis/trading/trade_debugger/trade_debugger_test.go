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
	t.ProxiedMainnetUser.AddBearerToken(t.Tc.ProductionLocalTemporalBearerToken)
	uni := web3_client.InitUniswapClient(ctx, t.ProxiedMainnetUser)
	uni.PrintOn = true
	uni.PrintLocal = true
	uni.Web3Client.IsAnvilNode = true
	uni.Web3Client.DurableExecution = false
	uni.DebugPrint = true

	at := artemis_realtime_trading.NewActiveTradingDebugger(&uni)
	td := NewTradeDebugger(at, t.MainnetWeb3User)
	t.Require().NotEmpty(td)
	t.td = td
}

func (t *ArtemisTradeDebuggerTestSuite) TestDebuggerInitEnv() {
	txHash := "0x5327295e1ed6d59faaf98d04697b0316fb8ad4b767d2e7f5addb3981c3b5d3b7"
	mevTx, err := t.td.lookupMevTx(ctx, txHash)
	t.Require().Nil(err)
	t.Require().NotEmpty(mevTx)
	tf := mevTx.TradePrediction
	err = t.td.ResetAndSetupPreconditions(ctx, tf)
	fmt.Println(tf.FrontRunTrade.AmountIn)
	fmt.Println(tf.FrontRunTrade.AmountOut)
	b, terr := t.td.dat.SimW3c().ReadERC20TokenBalance(ctx, tf.FrontRunTrade.AmountInAddr.String(), t.td.dat.SimW3c().PublicKey())
	t.Require().Nil(terr)
	t.Require().Equal(tf.FrontRunTrade.AmountIn.String(), b.String())
}
func (t *ArtemisTradeDebuggerTestSuite) TestReplayDebugger1() {
	txHash := "0x5327295e1ed6d59faaf98d04697b0316fb8ad4b767d2e7f5addb3981c3b5d3b7"
	h, err := t.td.lookupMevTx(ctx, txHash)
	t.Require().Nil(err)
	t.Require().NotEmpty(h)
	fmt.Println(h.HistoricalAnalysis.TradeMethod)
	fmt.Println(h.HistoricalAnalysis.EndReason)
	tf, serr := h.BinarySearch()
	t.Require().Nil(serr)
	t.Require().NotEmpty(tf)
}

func TestArtemisTradeDebuggerTestSuite(t *testing.T) {
	suite.Run(t, new(ArtemisTradeDebuggerTestSuite))
}
