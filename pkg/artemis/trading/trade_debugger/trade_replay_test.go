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

// 0x80ae3cc1748c10f42e591783001817b8a56b188eb1867282e396a8d99d583d00

func (t *ArtemisTradeDebuggerTestSuite) TestReplayer() {
	txHash := "0xb2fc21dfc699958484cf9b0c97e9afccb5e8274daabb30164f0f9fe0dffb82fc"

	err := t.td.Replay(ctx, txHash)
	t.NoError(err)

}
