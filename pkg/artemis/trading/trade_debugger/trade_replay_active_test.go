package artemis_trade_debugger

import artemis_reporting "github.com/zeus-fyi/olympus/pkg/artemis/trading/reporting"

func (t *ArtemisTradeDebuggerTestSuite) TestActiveReplay() {
	bg, err := artemis_reporting.GetBundlesProfitHistory(ctx, 0, 1)
	t.Assert().Nil(err)
	t.Assert().NotNil(bg)

	// TODO
}
