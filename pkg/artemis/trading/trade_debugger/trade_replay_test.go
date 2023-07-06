package artemis_trade_debugger

import "fmt"

func (t *ArtemisTradeDebuggerTestSuite) TestReplayDebugger() {
	txHash := "0x5327295e1ed6d59faaf98d04697b0316fb8ad4b767d2e7f5addb3981c3b5d3b7"

	hist, err := t.td.lookupMevTx(ctx, txHash)
	t.Require().Nil(err)
	t.Require().NotEmpty(hist)

	for _, h := range hist {
		fmt.Println(h.HistoricalAnalysis.TradeMethod)
	}
}

func (t *ArtemisTradeDebuggerTestSuite) TestDebugger() {
	txHash := "0x5327295e1ed6d59faaf98d04697b0316fb8ad4b767d2e7f5addb3981c3b5d3b7"
	_, err := t.td.getTxFromHash(ctx, txHash)
	t.Require().Nil(err)

	rx, err := t.td.getRxFromHash(ctx, txHash)
	t.Require().Nil(err)
	t.Require().NotNil(rx)
	fmt.Println(rx.BlockNumber.String())
}
