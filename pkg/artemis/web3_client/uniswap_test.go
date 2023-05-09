package web3_client

import "fmt"

func (s *Web3ClientTestSuite) TestUniswapMempoolFilter() {
	uni := InitUniswapV2Client(ctx)
	txMap, err := s.MainnetWeb3User.GetFilteredPendingMempoolTxs(ctx, uni.MevSmartContractTxMap)
	s.Require().Nil(err)
	s.Assert().NotEmpty(txMap)
	uni.MevSmartContractTxMap = txMap
	uni.ProcessTxs()
	count := len(uni.SwapExactTokensForTokensParamsSlice)
	count += len(uni.SwapTokensForExactTokensParamsSlice)
	count += len(uni.SwapExactETHForTokensParamsSlice)
	count += len(uni.SwapTokensForExactETHParamsSlice)
	count += len(uni.SwapExactTokensForETHParamsSlice)
	count += len(uni.SwapETHForExactTokensParamsSlice)
	fmt.Println("Total trades found", count)
}
