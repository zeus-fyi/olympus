package web3_client

//
//func (s *Web3ClientTestSuite) TestRawMempoolTxFilter() {
//	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
//	ForceDirToTestDirLocation()
//	s.LocalMainnetWeb3User.Web3Actions.Dial()
//	defer s.LocalMainnetWeb3User.Close()
//	mempool, err := s.LocalMainnetWeb3User.Web3Actions.GetTxPoolContent(ctx)
//	s.Require().NoError(err)
//	uni := InitUniswapClient(ctx, s.MainnetWeb3User)
//	uni.PrintOn = true
//	uni.PrintLocal = true
//	uni.Path = filepaths.Path{
//		PackageName: "",
//		DirIn:       "",
//		DirOut:      "./trade_analysis",
//		FnIn:        "",
//		FnOut:       "",
//		Env:         "",
//	}
//	txMap, err := ProcessMempoolTxs(ctx, mempool["mempool"], uni.MevSmartContractTxMap)
//	s.Require().Nil(err)
//	s.Assert().NotEmpty(txMap)
//	uni.MevSmartContractTxMap = txMap
//	uni.ProcessTxs(ctx)
//	count := len(uni.SwapExactTokensForTokensParamsSlice)
//	fmt.Println("Total SwapExactTokensForTokensParamsSlice found", len(uni.SwapExactTokensForTokensParamsSlice))
//	count += len(uni.SwapTokensForExactTokensParamsSlice)
//	fmt.Println("Total SwapTokensForExactTokensParamsSlice found", len(uni.SwapTokensForExactTokensParamsSlice))
//	count += len(uni.SwapExactETHForTokensParamsSlice)
//	fmt.Println("Total SwapExactETHForTokensParamsSlice found", len(uni.SwapExactETHForTokensParamsSlice))
//	count += len(uni.SwapTokensForExactETHParamsSlice)
//	fmt.Println("Total SwapTokensForExactETHParamsSlice found", len(uni.SwapTokensForExactETHParamsSlice))
//	count += len(uni.SwapExactTokensForETHParamsSlice)
//	fmt.Println("Total SwapExactTokensForETHParamsSlice found", len(uni.SwapExactTokensForETHParamsSlice))
//	count += len(uni.SwapETHForExactTokensParamsSlice)
//	fmt.Println("Total SwapETHForExactTokensParamsSlice found", len(uni.SwapETHForExactTokensParamsSlice))
//	fmt.Println("Total trades found", count)
//}
