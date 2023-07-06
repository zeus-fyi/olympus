package artemis_trade_debugger

import (
	"fmt"
)

/*
type TradeExecutionFlow struct {
	CurrentBlockNumber *big.Int                           `json:"currentBlockNumber"`
	Tx                 *types.Transaction                 `json:"tx"`
	Trade              Trade                              `json:"trade"`
	InitialPair        *uniswap_pricing.UniswapV2Pair     `json:"initialPair,omitempty"`
	InitialPairV3      *uniswap_pricing.UniswapV3Pair     `json:"initialPairV3,omitempty"`
	FrontRunTrade      artemis_trading_types.TradeOutcome `json:"frontRunTrade"`
	UserTrade          artemis_trading_types.TradeOutcome `json:"userTrade"`
	SandwichTrade      artemis_trading_types.TradeOutcome `json:"sandwichTrade"`
	SandwichPrediction SandwichTradePrediction            `json:"sandwichPrediction"`
}
*/

func (t *ArtemisTradeDebuggerTestSuite) TestReplayDebugger() {
	txHash := "0x5327295e1ed6d59faaf98d04697b0316fb8ad4b767d2e7f5addb3981c3b5d3b7"
	hist, err := t.td.lookupMevTx(ctx, txHash)
	t.Require().Nil(err)
	t.Require().NotEmpty(hist)

	for _, h := range hist {
		fmt.Println(h.HistoricalAnalysis.TradeMethod)
		fmt.Println(h.HistoricalAnalysis.EndReason)
		fmt.Println(h.TradeParams)

		tf, serr := h.BinarySearch()
		t.Require().Nil(serr)
		fmt.Println(tf)
	}
	t.Assert().NotEmpty(hist)

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
