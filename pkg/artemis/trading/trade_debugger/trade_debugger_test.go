package artemis_trade_debugger

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_realtime_trading "github.com/zeus-fyi/olympus/pkg/artemis/trading"
	artemis_trading_auxiliary "github.com/zeus-fyi/olympus/pkg/artemis/trading/auxiliary"
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
	uni.PrintOn = true
	uni.PrintLocal = false
	uni.DebugPrint = true
	uni.Web3Client.IsAnvilNode = true
	uni.Web3Client.DurableExecution = false
	a := artemis_trading_auxiliary.AuxiliaryTradingUtils{
		U: &uni,
	}
	at := artemis_realtime_trading.NewActiveTradingModuleWithoutMetrics(&a)
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
	b, terr := t.td.UniswapClient.Web3Client.ReadERC20TokenBalance(ctx, tf.FrontRunTrade.AmountInAddr.String(), t.td.UniswapClient.Web3Client.PublicKey())
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
