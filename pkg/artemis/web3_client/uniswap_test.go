package web3_client

import "fmt"

func (s *Web3ClientTestSuite) TestUniswapMempoolFilter() {
	uni := InitUniswapV2Client(ctx)
	txMap, err := s.MainnetWeb3User.GetFilteredPendingMempoolTxs(ctx, uni.MevSmartContractTxMap)
	s.Require().Nil(err)
	s.Assert().NotEmpty(txMap)
	uni.MevSmartContractTxMap = txMap
	uni.ProcessTxs()

	fmt.Println("SwapExactTokensForTokensParamsSlice", uni.SwapExactTokensForTokensParamsSlice)
	fmt.Println("SwapTokensForExactTokensParamsSlice", uni.SwapTokensForExactTokensParamsSlice)
	fmt.Println("SwapExactETHForTokensParamsSlice", uni.SwapExactETHForTokensParamsSlice)
	fmt.Println("SwapTokensForExactETHParamsSlice", uni.SwapTokensForExactETHParamsSlice)
	fmt.Println("SwapExactTokensForETHParamsSlice", uni.SwapExactTokensForETHParamsSlice)
	fmt.Println("SwapETHForExactTokensParamsSlice", uni.SwapETHForExactTokensParamsSlice)

	count := len(uni.SwapExactTokensForTokensParamsSlice)
	count += len(uni.SwapTokensForExactTokensParamsSlice)
	count += len(uni.SwapExactETHForTokensParamsSlice)
	count += len(uni.SwapTokensForExactETHParamsSlice)
	count += len(uni.SwapExactTokensForETHParamsSlice)
	count += len(uni.SwapETHForExactTokensParamsSlice)

}
