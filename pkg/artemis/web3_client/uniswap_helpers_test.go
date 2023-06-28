package web3_client

func (s *Web3ClientTestSuite) TestUniswapSortTokens() {
	p := UniswapV2Pair{}
	err := p.PairForV2(PepeContractAddr, WETH9ContractAddress)
	s.Require().Nil(err)
	s.Require().Equal("0xA43fe16908251ee70EF74718545e4FE6C5cCEc9f", p.PairContractAddr)
	s.Require().Equal(p.Token0.String(), PepeContractAddr)
	s.Require().Equal(p.Token1.String(), WETH9ContractAddress)
}
