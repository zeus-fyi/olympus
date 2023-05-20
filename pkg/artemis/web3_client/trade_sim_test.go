package web3_client

/*
	modEven :=  tokenSellAmount.Mod(tokenSellAmount, big.NewInt(2))
	if modEven.String() == "0" {
		tokenSellAmount = tokenSellAmount.Div(tokenSellAmount, big.NewInt(2))
	} else {
		tokenSellAmount = tokenSellAmount.Add(tokenSellAmount, big.NewInt(1))
		tokenSellAmount = tokenSellAmount.Div(tokenSellAmount, big.NewInt(2))
	}
*/

func (s *Web3ClientTestSuite) TestJson() {
	//amountOutMin, _ := new(big.Int).SetString("746627207819418433569734379647", 10)
	//te := TradeExecutionFlow{
	//	CurrentBlockNumber: nil,
	//	Tx:                 nil,
	//	Trade: Trade{
	//		TradeMethod:                    "swapExactETHForTokens",
	//		SwapETHForExactTokensParams:    nil,
	//		SwapTokensForExactTokensParams: nil,
	//		SwapExactTokensForTokensParams: nil,
	//		SwapExactETHForTokensParams: &SwapExactETHForTokensParams{
	//			AmountOutMin: amountOutMin,
	//		},
	//		SwapExactTokensForETHParams: nil,
	//		SwapTokensForExactETHParams: nil,
	//	},
	//	InitialPair:        UniswapV2Pair{},
	//	FrontRunTrade:      TradeOutcome{},
	//	SandwichTrade:      TradeOutcome{},
	//	SandwichPrediction: SandwichTradePrediction{},
	//}
	//
	//b, _ := json.Marshal(te)
	//fmt.Println(string(b))
}

func (s *Web3ClientTestSuite) TestTradeSim() {
	//apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	//ForceDirToTestDirLocation()
	//uni := InitUniswapV2Client(ctx, s.MainnetWeb3User)
	//mevTxs, err := artemis_validator_service_groups_models.SelectMempoolTxAtBlockNumber(ctx, hestia_req_types.EthereumMainnetProtocolNetworkID, 17290284)
	//s.Require().Nil(err)
	//s.Require().NotEmpty(mevTxs)
	//
	//for _, mevTx := range mevTxs {
	//	tf := TradeExecutionFlow{}
	//	b := []byte(mevTx.TxFlowPrediction)
	//	berr := json.Unmarshal(b, &tf)
	//	s.Require().Nil(berr)
	//
	//	if tf.FrontRunTrade.AmountIn == nil {
	//		continue
	//	}
	//	fmt.Println(tf.Trade.TradeMethod)
	//	b, berr = json.MarshalIndent(tf, "", "  ")
	//	s.Require().Nil(berr)
	//	//fmt.Println(string(b))
	//	executedProfit := uni.TradeSimStep(tf)
	//	fmt.Println("binary search sell amount", tf.SandwichPrediction.SellAmount.String())
	//	fmt.Println("binary search max profit", tf.SandwichPrediction.ExpectedProfit.String())
	//	fmt.Println("profit from execution path", executedProfit.String())
	//	fmt.Println("profit token type", tf.SandwichTrade.AmountOutAddr.String())
	//	s.Assert().Equal(tf.SandwichPrediction.ExpectedProfit.String(), executedProfit.String())
	//
	//	sellAmount, maxProfit := uni.TradeSim(tf)
	//	fmt.Println("linear search sell amount", sellAmount.String())
	//	fmt.Println("linear search max profit", maxProfit.String())
	//	tf.FrontRunTrade.AmountIn = sellAmount
	//}
}
