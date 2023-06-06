package web3_client

func (s *Web3ClientTestSuite) TestMevTxFilter() {
	uni := InitUniswapClient(ctx, s.MainnetWeb3User)
	txMap, err := s.MainnetWeb3User.GetFilteredPendingMempoolTxs(ctx, uni.MevSmartContractTxMap)
	s.Require().Nil(err)
	s.Assert().NotEmpty(txMap)
	uni.MevSmartContractTxMap = txMap
	uni.ProcessTxs(ctx)
}
