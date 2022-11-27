package web3_client

func (s *Web3ClientTestSuite) TestWeb3SendEther() {
	b, err := s.GoerliWeb3User.GetCurrentBalance(ctx)

	s.Require().Nil(err)
	s.Assert().NotNil(b)

	send, err := s.GoerliWeb3User.Send(ctx, s.GoerliWeb3User2.Address(), Finney, nil, 0)
	s.Require().Nil(err)
	s.Assert().NotNil(send)
}
