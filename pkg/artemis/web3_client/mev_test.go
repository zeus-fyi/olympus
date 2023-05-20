package web3_client

func (s *Web3ClientTestSuite) TestRawMempoolTxFilter() {
	uni := InitUniswapV2Client(ctx, s.LocalMainnetWeb3User)
	s.LocalMainnetWeb3User.Web3Actions.Dial()
	defer s.LocalMainnetWeb3User.Close()
	mempool, err := s.LocalMainnetWeb3User.Web3Actions.GetTxPoolContent(ctx)
	s.Require().NoError(err)
	txMap, err := s.LocalMainnetWeb3User.ProcessUnvalidatedMempoolTxs(ctx, mempool, uni.MevSmartContractTxMap)
	s.Require().Nil(err)
	s.Assert().NotEmpty(txMap)
}
func (s *Web3ClientTestSuite) TestMevTxFilter() {
	uni := InitUniswapV2Client(ctx, s.MainnetWeb3User)
	txMap, err := s.MainnetWeb3User.GetFilteredPendingMempoolTxs(ctx, uni.MevSmartContractTxMap)
	s.Require().Nil(err)
	s.Assert().NotEmpty(txMap)
	uni.MevSmartContractTxMap = txMap
	uni.ProcessTxs(ctx)
}
