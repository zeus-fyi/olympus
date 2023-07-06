package artemis_trade_debugger

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

func (t *ArtemisTradeDebuggerTestSuite) TestReplayer() {
	txHash := "0x5327295e1ed6d59faaf98d04697b0316fb8ad4b767d2e7f5addb3981c3b5d3b7"

	err := t.td.Replay(ctx, txHash)
	t.NoError(err)

}
