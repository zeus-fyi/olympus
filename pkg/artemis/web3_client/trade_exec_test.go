package web3_client

func (s *Web3ClientTestSuite) TestTradeExec() {
	forceDirToLocation()
	swapAbi, bc, err := LoadSwapAbiPayload()
	s.Require().NoError(err)
	s.Require().NotNil(swapAbi)
	s.Require().NotNil(bc)

}
